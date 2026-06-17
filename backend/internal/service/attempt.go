package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"test-constructor/internal/client"
	"test-constructor/internal/domain"
	"test-constructor/internal/dto"
	"test-constructor/internal/repository"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AttemptService interface {
	StartAttempt(userID uint, link string, req dto.StartAttemptRequest) (*dto.StartAttemptResponse, error)
	GetActiveAttempt(userID uint) (*dto.StartAttemptResponse, error)
	FinishAttempt(userID uint, req dto.FinishAttemptRequest) (*dto.FinishAttemptResponse, error)
}

type attemptService struct {
	attemptRepo        repository.AttemptRepository
	answerRepo         repository.AnswerRepository
	eventConfigRepo    repository.EventConfigRepository
	extraThresholdRepo repository.ExtraThresholdRepository
	questionRepo       repository.QuestionRepository
	userEventRepo      repository.UserEventRepository
	txManager          repository.TransactionManager
	crmClient          client.CRMClient
}

func NewAttemptService(
	attemptRepo repository.AttemptRepository,
	answerRepo repository.AnswerRepository,
	eventConfigRepo repository.EventConfigRepository,
	extraThresholdRepo repository.ExtraThresholdRepository,
	questionRepo repository.QuestionRepository,
	userEventRepo repository.UserEventRepository,
	txManager repository.TransactionManager,
	crmClient client.CRMClient,
) AttemptService {
	return &attemptService{
		attemptRepo:        attemptRepo,
		answerRepo:         answerRepo,
		eventConfigRepo:    eventConfigRepo,
		extraThresholdRepo: extraThresholdRepo,
		questionRepo:       questionRepo,
		userEventRepo:      userEventRepo,
		txManager:          txManager,
		crmClient:          crmClient,
	}
}

func (s *attemptService) StartAttempt(userID uint, link string, req dto.StartAttemptRequest) (*dto.StartAttemptResponse, error) {
	config, err := s.eventConfigRepo.FindByTestLink(link)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("тест не найден")
		}
		return nil, fmt.Errorf("ошибка поиска теста: %w", err)
	}

	_, err = s.attemptRepo.FindActiveByUser(userID)
	if err == nil {
		return nil, errors.New("у вас уже есть активная попытка. Завершите её перед началом новой")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("ошибка проверки активных попыток: %w", err)
	}

	existingAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, config.ConfigID)
	if err == nil && existingAttempt.EndTime != nil {
		return nil, errors.New("вы уже прошли этот тест")
	}

	if config.IsExtra {
		if !s.hasAccessToExtraConfig(userID, config.ConfigID) {
			return nil, errors.New("у вас нет доступа к этому дополнительному тесту")
		}
	}

	now := time.Now()
	maxScore, err := s.questionRepo.GetMaxScoreByTestID(config.TestID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения максимального балла: %w", err)
	}

	attempt := domain.Attempt{
		ConfigID:      config.ConfigID,
		ApplicationID: req.ApplicationID,
		InternID:      userID,
		StartTime:     now,
		MaxScore:      maxScore,
	}
	if err = s.attemptRepo.Create(&attempt); err != nil {
		return nil, fmt.Errorf("ошибка создания попытки: %w", err)
	}

	expiresAt := now.Add(time.Duration(config.TimeLimit) * time.Minute)
	err = s.crmClient.CreateTestSession(req.ApplicationID, config.TestID, fmt.Sprintf("%d", attempt.AttemptID), expiresAt)
	if err != nil {
		fmt.Printf("Ошибка создания сессии в CRM: %v\n", err)
	}

	publicQuestions, err := s.preparePublicQuestions(config.TestID)
	if err != nil {
		return nil, err
	}

	return &dto.StartAttemptResponse{
		AttemptID:   attempt.AttemptID,
		ConfigID:    config.ConfigID,
		TestID:      config.TestID,
		Title:       config.Test.Title,
		Description: config.Test.Description,
		TimeLimit:   config.TimeLimit,
		Threshold:   config.Threshold,
		Questions:   publicQuestions,
	}, nil
}

func (s *attemptService) GetActiveAttempt(userID uint) (*dto.StartAttemptResponse, error) {
	attempt, err := s.attemptRepo.FindActiveByUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("активная попытка не найдена")
		}
		return nil, fmt.Errorf("ошибка поиска активной попытки: %w", err)
	}

	config, err := s.eventConfigRepo.FindByID(attempt.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	publicQuestions, err := s.preparePublicQuestions(config.TestID)
	if err != nil {
		return nil, err
	}

	return &dto.StartAttemptResponse{
		AttemptID:   attempt.AttemptID,
		ConfigID:    config.ConfigID,
		TestID:      config.TestID,
		Title:       config.Test.Title,
		Description: config.Test.Description,
		TimeLimit:   config.TimeLimit,
		Threshold:   config.Threshold,
		Questions:   publicQuestions,
	}, nil
}

func (s *attemptService) FinishAttempt(userID uint, req dto.FinishAttemptRequest) (*dto.FinishAttemptResponse, error) {
	attempt, err := s.attemptRepo.FindActiveByUser(userID)
	if err != nil {
		return nil, errors.New("активная попытка не найдена")
	}

	config, err := s.eventConfigRepo.FindByID(attempt.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}
	questions, err := s.questionRepo.FindByTestID(config.TestID)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки вопросов: %w", err)
	}

	questionMap := make(map[uint]domain.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	userPoints := 0
	answers := make([]domain.Answer, 0, len(req.UserAnswers))
	for _, ansInfo := range req.UserAnswers {
		question, ok := questionMap[ansInfo.QuestionID]
		if !ok {
			return nil, fmt.Errorf("вопрос с ID %d не найден", ansInfo.QuestionID)
		}
		if question.TestID != config.TestID {
			return nil, errors.New("ответы не соответствуют тесту")
		}

		correct := checkAnswer(question.Type, question.Options, ansInfo.Answer)
		points := 0
		if correct {
			points = question.Points
			userPoints += points
		}

		answerJSON, _ := json.Marshal(ansInfo.Answer)
		answers = append(answers, domain.Answer{
			AttemptID:    attempt.AttemptID,
			QuestionID:   question.ID,
			InternAnswer: datatypes.JSON(answerJSON),
			IsCorrect:    correct,
			Points:       points,
		})
	}

	tx, err := s.txManager.Begin()
	if err != nil {
		return nil, errors.New("ошибка начала транзакции")
	}

	if err = s.answerRepo.CreateBatchWithTx(tx, answers); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("ошибка сохранения ответов: %w", err)
	}

	now := time.Now()
	attempt.EndTime = &now
	attempt.Score = userPoints
	attempt.Passed = userPoints >= config.Threshold

	if err = s.attemptRepo.UpdateWithTx(tx, attempt); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("ошибка обновления попытки: %w", err)
	}

	if err = s.txManager.Commit(tx); err != nil {
		return nil, errors.New("ошибка сохранения результатов")
	}

	allCompleted := s.checkAllTestsCompleted(userID, config.EventID, config.SpecializationID)

	eventPassed := s.checkEventPassed(userID, config.EventID, config.SpecializationID)

	if allCompleted || !s.hasAvailableReplacements(userID, config.EventID, config.SpecializationID) {
		crmResult := client.CRMResultData{
			SessionID:   fmt.Sprintf("%d", attempt.AttemptID),
			TestID:      fmt.Sprintf("%d", config.TestID),
			Score:       userPoints,
			MaxScore:    attempt.MaxScore,
			IsPassed:    eventPassed,
			CompletedAt: now.Format("2006-01-02T15:04:05Z"),
			StartedAt:   attempt.StartTime.Format("2006-01-02T15:04:05Z"),
		}
		if err := s.crmClient.SendTestResult(attempt.ApplicationID, crmResult); err != nil {
			fmt.Printf("Ошибка отправки результатов в CRM: %v\n", err)
		}
	}

	resultText := config.FailText
	if attempt.Passed {
		resultText = config.SuccessText
	}

	return &dto.FinishAttemptResponse{
		Result:        resultText,
		Score:         userPoints,
		MaxTestPoints: attempt.MaxScore,
		Passed:        attempt.Passed,
		AllCompleted:  allCompleted,
	}, nil
}

func checkAnswer(qType domain.QType, optionsJSON datatypes.JSON, answer dto.UserAnswer) bool {
	var opts domain.QuestionOptions
	if err := json.Unmarshal(optionsJSON, &opts); err != nil {
		return false
	}

	switch qType {
	case domain.SingleChoice, domain.MultipleChoice:
		if len(answer.Choices) != len(opts.Choices) {
			return false
		}
		for i, c := range opts.Choices {
			if c.IsTrue != answer.Choices[i] {
				return false
			}
		}
		return true
	case domain.Matching:
		if len(answer.MatchingPairs) != len(opts.MatchingPairs) {
			return false
		}
		for i, p := range opts.MatchingPairs {
			if p.LeftColumn != answer.MatchingPairs[i].Left || p.RightColumn != answer.MatchingPairs[i].Right {
				return false
			}
		}
		return true
	case domain.CorrectOrder:
		if len(answer.Sequence) != len(opts.Sequence) {
			return false
		}
		for i, s := range opts.Sequence {
			if s.Order != answer.Sequence[i].Order {
				return false
			}
		}
		return true
	case domain.TextInput:
		for _, correct := range opts.CorrectInput {
			if (opts.CaseSensitive && answer.UserInput == correct) ||
				(!opts.CaseSensitive && strings.EqualFold(answer.UserInput, correct)) {
				return true
			}
		}
		return false
	}
	return false
}

func (s *attemptService) hasAccessToExtraConfig(userID, extraConfigID uint) bool {
	threshold, err := s.extraThresholdRepo.FindByExtraConfigID(extraConfigID)
	if err != nil {
		return false
	}

	attempt, err := s.attemptRepo.FindByUserAndConfig(userID, threshold.ConfigID)
	if err != nil || attempt.EndTime == nil {
		return false
	}
	return attempt.Score >= threshold.Threshold
}

func (s *attemptService) checkAllTestsCompleted(userID, eventID, specializationID uint) bool {
	mainConfigs, err := s.eventConfigRepo.FindMainConfigsByEventAndSpec(eventID, specializationID)
	if err != nil {
		return false
	}

	for _, mainCfg := range mainConfigs {
		mainAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, mainCfg.ConfigID)
		if err != nil || mainAttempt.EndTime == nil {
			return false
		}
		if mainAttempt.Passed {
			continue
		}

		if !s.hasCompletedReplacement(userID, mainCfg.ConfigID) {
			return false
		}
	}
	return true
}

func (s *attemptService) checkEventPassed(userID, eventID, specializationID uint) bool {
	mainConfigs, err := s.eventConfigRepo.FindMainConfigsByEventAndSpec(eventID, specializationID)
	if err != nil || len(mainConfigs) == 0 {
		return false
	}

	for _, mainCfg := range mainConfigs {
		mainAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, mainCfg.ConfigID)
		if err != nil || mainAttempt.EndTime == nil {
			return false
		}

		if mainAttempt.Passed {
			continue
		}

		if !s.hasPassedReplacement(userID, mainCfg.ConfigID) {
			return false
		}
	}
	return true
}

func (s *attemptService) hasPassedReplacement(userID, mainConfigID uint) bool {
	replacements, _ := s.extraThresholdRepo.FindReplacementsForConfigID(mainConfigID)
	mainAttempt, _ := s.attemptRepo.FindByUserAndConfig(userID, mainConfigID)

	for _, replacement := range replacements {
		if mainAttempt.Score < replacement.Threshold {
			continue
		}
		extraAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, replacement.ExtraConfigID)
		if err == nil && extraAttempt.EndTime != nil && extraAttempt.Passed {
			return true
		}
	}
	return false
}

func (s *attemptService) hasCompletedReplacement(userID, mainConfigID uint) bool {
	replacements, _ := s.extraThresholdRepo.FindReplacementsForConfigID(mainConfigID)
	mainAttempt, _ := s.attemptRepo.FindByUserAndConfig(userID, mainConfigID)

	for _, replacement := range replacements {
		if mainAttempt.Score < replacement.Threshold {
			continue
		}
		extraAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, replacement.ExtraConfigID)
		if err == nil && extraAttempt.EndTime != nil {
			return true
		}
	}
	return false
}

func (s *attemptService) hasAvailableReplacements(userID, eventID, specializationID uint) bool {
	mainConfigs, _ := s.eventConfigRepo.FindMainConfigsByEventAndSpec(eventID, specializationID)

	for _, mainCfg := range mainConfigs {
		mainAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, mainCfg.ConfigID)
		if err != nil || mainAttempt.EndTime == nil {
			continue
		}

		if mainAttempt.Passed {
			continue
		}

		replacements, _ := s.extraThresholdRepo.FindReplacementsForConfigID(mainCfg.ConfigID)
		for _, replacement := range replacements {
			if mainAttempt.Score >= replacement.Threshold {
				extraAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, replacement.ExtraConfigID)
				if err != nil || extraAttempt.EndTime == nil {
					return true
				}
			}
		}
	}
	return false
}

func (s *attemptService) preparePublicQuestions(testID uint) ([]dto.QuestionPublic, error) {
	questions, err := s.questionRepo.FindByTestID(testID)
	if err != nil {
		return nil, err
	}

	public := make([]dto.QuestionPublic, len(questions))
	for i, q := range questions {
		var opts domain.QuestionOptions
		json.Unmarshal(q.Options, &opts)

		pub := dto.QuestionPublic{
			QuestionID:  q.ID,
			Text:        q.Text,
			Points:      q.Points,
			OrderNumber: q.OrderNumber,
			Type:        string(q.Type),
		}

		switch q.Type {
		case domain.SingleChoice, domain.MultipleChoice:
			choices := make([]string, len(opts.Choices))
			for j, c := range opts.Choices {
				choices[j] = c.Text
			}

			rand.Shuffle(len(choices), func(i, j int) {
				choices[i], choices[j] = choices[j], choices[i]
			})
			pub.Options.Choices = choices
		case domain.Matching:
			left := make([]string, len(opts.MatchingPairs))
			right := make([]string, len(opts.MatchingPairs))
			for j, p := range opts.MatchingPairs {
				left[j] = p.LeftColumn
				right[j] = p.RightColumn
			}
			rand.Shuffle(len(right), func(i, j int) {
				right[i], right[j] = right[j], right[i]
			})
			pub.Options.Matching = &dto.PublicMatching{LeftColumn: left, RightColumn: right}
		case domain.CorrectOrder:
			texts := make([]string, len(opts.Sequence))
			for j, s := range opts.Sequence {
				texts[j] = s.Text
			}
			rand.Shuffle(len(texts), func(i, j int) {
				texts[i], texts[j] = texts[j], texts[i]
			})
			pub.Options.Sequence = texts
		case domain.TextInput:
			pub.Options.CaseSensitive = opts.CaseSensitive
		}
		public[i] = pub
	}
	return public, nil
}
