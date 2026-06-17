package service

import (
	"fmt"
	"sort"
	"test-constructor/internal/client"
	"test-constructor/internal/domain"
	"test-constructor/internal/dto"
	"test-constructor/internal/repository"
)

type StatisticsService interface {
	GetInternList() (*dto.GetUsersResponse, error)
	GetUserStatistics(userID uint) (*dto.UserStatisticsResponse, error)
	GetEventStatistics(eventID uint, filter *dto.StatisticsFilter) (*dto.StatisticsResponse, error)
}

type statisticsService struct {
	statsRepo repository.StatisticsRepository
	crmClient client.CRMClient
}

func NewStatisticsService(
	statsRepo repository.StatisticsRepository,
	crmClient client.CRMClient,
) StatisticsService {
	return &statisticsService{
		statsRepo: statsRepo,
		crmClient: crmClient,
	}
}

func (s *statisticsService) GetInternList() (*dto.GetUsersResponse, error) {
	users, err := s.statsRepo.FindInternsByRole("intern")
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка стажёров: %w", err)
	}

	interns := make([]dto.UserInfo, len(users))
	for i, user := range users {
		interns[i] = dto.UserInfo{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
		}
	}

	return &dto.GetUsersResponse{
		Users: interns,
	}, nil
}

func (s *statisticsService) GetUserStatistics(userID uint) (*dto.UserStatisticsResponse, error) {
	user, err := s.statsRepo.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}

	attempts, err := s.statsRepo.FindCompletedAttemptsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения попыток: %w", err)
	}

	events, err := s.crmClient.GetEvents()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения мероприятий: %w", err)
	}

	eventMap := make(map[uint]string)
	for _, event := range events {
		eventMap[uint(event.ID)] = event.Name
	}

	attemptDetails := make([]dto.UserAttemptDetail, 0, len(attempts))
	for _, attempt := range attempts {
		cfg := attempt.EventConfig

		eventName, exists := eventMap[cfg.EventID]
		if !exists {
			eventName = fmt.Sprintf("Event #%d", cfg.EventID)
		}

		maxScore := 0
		questionMap := make(map[uint]domain.Question)
		for _, q := range cfg.Test.Questions {
			maxScore += q.Points
			questionMap[q.ID] = q
		}

		questions := s.getQuestionsStat(attempt.Answers, questionMap)

		detail := dto.UserAttemptDetail{
			AttemptID: attempt.AttemptID,
			TestTitle: cfg.Test.Title,
			EventName: eventName,
			IsExtra:   cfg.IsExtra,
			Score:     attempt.Score,
			MaxScore:  maxScore,
			Passed:    attempt.Passed,
			Questions: questions,
		}

		attemptDetails = append(attemptDetails, detail)
	}

	return &dto.UserStatisticsResponse{
		UserID:    user.ID,
		FirstName: user.Name,
		LastName:  user.Surname,
		Email:     user.Email,
		Attempts:  attemptDetails,
	}, nil
}

func (s *statisticsService) GetEventStatistics(eventID uint, filter *dto.StatisticsFilter) (*dto.StatisticsResponse, error) {
	var isExtra *bool
	if filter != nil {
		isExtra = filter.IsExtra
	}

	configs, err := s.statsRepo.FindConfigsByEventID(eventID, isExtra)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения конфигураций: %w", err)
	}

	if len(configs) == 0 {
		return &dto.StatisticsResponse{
			Attempts: []dto.UserAttemptInfo{},
		}, nil
	}

	configIDs := make([]uint, len(configs))
	configMap := make(map[uint]domain.EventConfig)
	questionMap := make(map[uint]map[uint]domain.Question)

	for i, cfg := range configs {
		configIDs[i] = cfg.ConfigID
		configMap[cfg.ConfigID] = cfg

		questionMap[cfg.ConfigID] = make(map[uint]domain.Question)
		for _, q := range cfg.Test.Questions {
			questionMap[cfg.ConfigID][q.ID] = q
		}
	}

	attempts, err := s.statsRepo.FindCompletedAttemptsByConfigIDs(configIDs)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения попыток: %w", err)
	}

	userAttempts := make([]dto.UserAttemptInfo, 0, len(attempts))
	for _, attempt := range attempts {
		cfg, cfgExists := configMap[attempt.ConfigID]
		if !cfgExists {
			continue
		}

		maxScore := 0
		for _, question := range cfg.Test.Questions {
			maxScore += question.Points
		}

		timeSpent := 0
		if attempt.EndTime != nil {
			duration := attempt.EndTime.Sub(attempt.StartTime)
			timeSpent = int(duration.Minutes())
		}

		questions := s.getQuestionsStat(attempt.Answers, questionMap[cfg.ConfigID])

		attemptInfo := dto.UserAttemptInfo{
			UserID:    attempt.InternID,
			FirstName: attempt.User.Name,
			LastName:  attempt.User.Surname,
			Email:     attempt.User.Email,
			Score:     attempt.Score,
			MaxScore:  maxScore,
			Passed:    attempt.Passed,
			IsExtra:   attempt.EventConfig.IsExtra,
			TimeSpent: timeSpent,
			Questions: questions,
		}

		userAttempts = append(userAttempts, attemptInfo)
	}

	return &dto.StatisticsResponse{
		Attempts: userAttempts,
	}, nil
}

func (s *statisticsService) getQuestionsStat(answers []domain.Answer, questionsMap map[uint]domain.Question) []dto.QuestionStatInfo {
	var questionsStat []dto.QuestionStatInfo

	for _, answer := range answers {
		question, exists := questionsMap[answer.QuestionID]
		if !exists {
			if answer.Question.ID != 0 {
				question = answer.Question
			} else {
				continue
			}
		}

		questionInfo := dto.QuestionStatInfo{
			Text:         question.Text,
			Points:       answer.Points,
			MaxPoints:    question.Points,
			IsCorrect:    answer.IsCorrect,
			QuestionType: string(question.Type),
			OrderNumber:  question.OrderNumber,
		}
		questionsStat = append(questionsStat, questionInfo)
	}

	sort.Slice(questionsStat, func(i, j int) bool {
		return questionsStat[i].OrderNumber < questionsStat[j].OrderNumber
	})

	return questionsStat
}
