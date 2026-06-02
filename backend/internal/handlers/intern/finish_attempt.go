package intern

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"test-constructor/config"
	"test-constructor/internal/auth"
	"test-constructor/internal/database"
	"test-constructor/internal/middleware"
	"test-constructor/internal/models"
	"time"

	"gorm.io/datatypes"
)

type CRMResultData struct {
	SessionID   string `json:"session_id"`
	TestID      string `json:"test_id"`
	Score       int    `json:"score"`
	MaxScore    int    `json:"max_score"`
	IsPassed    bool   `json:"is_passed"`
	CompletedAt string `json:"completed_at"`
	StartedAt   string `json:"started_at"`
}

type FinishAttemptRequest struct {
	UserAnswers []UserAnswerInfo
}

type UserAnswerInfo struct {
	QuestionID uint       `json:"question_id"`
	Answer     UserAnswer `json:"answer"`
}

type UserAnswer struct {
	Choices       []bool                `json:"choices,omitempty"`
	MatchingPairs []models.MatchingPair `json:"matching,omitempty"`
	UserInput     string                `json:"user_input,omitempty"`
	Sequence      []models.SequenceItem `json:"sequence,omitempty"`
}

type FinishAttemptResponse struct {
	Result        string `json:"result"`
	Score         int    `json:"score"`
	MaxTestPoints int    `json:"max_test_points"`
	Passed        bool   `json:"passed"`
	AllCompleted  bool   `json:"all_completed"` // Все ли тесты мероприятия пройдены
}

// @Summary Завершить тест
// @Security ApiKeyAuth
// @Description Получение ответов стажёра и проверка завершения всех тестов
// @Tags intern
// @Accept json
// @Produce json
// @Param answers body FinishAttemptRequest true "Answers object"
// @Success 201 {object} FinishAttemptResponse
// @Router /api/intern/attempt/finish [post]
func FinishAttempt(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var req FinishAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неправильный формат запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Находим активную попытку
	var attempt models.Attempt
	if err := database.DB.Preload("EventConfig").
		Preload("EventConfig.Test").
		Preload("EventConfig.ExtraThreshold").
		Where("intern_id = ? AND end_time IS NULL", claims.UserID).
		First(&attempt).Error; err != nil {
		http.Error(w, "Активная попытка не найдена", http.StatusNotFound)
		return
	}

	var questions []models.Question
	if err := database.DB.Where("test_id = ?", attempt.EventConfig.TestID).
		Find(&questions).Error; err != nil {
		http.Error(w, "Ошибка загрузки вопросов", http.StatusInternalServerError)
		return
	}

	questionMap := make(map[uint]models.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	userPoints := 0
	maxPoints := attempt.MaxScore
	correctAnswersCount := 0

	for _, answerInfo := range req.UserAnswers {
		question, exists := questionMap[answerInfo.QuestionID]
		if !exists {
			http.Error(w, fmt.Sprintf("Вопрос с ID %d не найден", answerInfo.QuestionID), http.StatusNotFound)
			return
		}

		if question.TestID != attempt.EventConfig.TestID {
			http.Error(w, "Ответы не соответствуют тесту", http.StatusBadRequest)
			return
		}

		answer := answerInfo.Answer
		var options models.QuestionOptions
		if err := json.Unmarshal(question.Options, &options); err != nil {
			http.Error(w, "Ошибка формата вопроса", http.StatusInternalServerError)
			return
		}

		correct := checkAnswer(question.Type, options, answer)

		if correct {
			userPoints += question.Points
			correctAnswersCount++
		}

		answerJSON, err := json.Marshal(answer)
		if err != nil {
			http.Error(w, "Ошибка сохранения ответа", http.StatusInternalServerError)
			return
		}

		userAnswer := models.Answer{
			QuestionID:   question.ID,
			AttemptID:    attempt.AttemptID,
			InternAnswer: datatypes.JSON(answerJSON),
			IsCorrect:    correct,
			Points: func() float64 {
				if correct {
					return float64(question.Points)
				}
				return 0
			}(),
		}

		if err := database.DB.Create(&userAnswer).Error; err != nil {
			http.Error(w, "Ошибка сохранения ответа в БД", http.StatusInternalServerError)
			return
		}
	}

	percentage := 0.0
	if maxPoints > 0 {
		percentage = float64(userPoints) / float64(maxPoints) * 100
	}

	passed := percentage >= float64(attempt.EventConfig.Threshold)

	now := time.Now()
	attempt.EndTime = &now
	attempt.Score = userPoints
	attempt.Passed = passed

	if err := database.DB.Save(&attempt).Error; err != nil {
		http.Error(w, "Ошибка сохранения попытки", http.StatusInternalServerError)
		return
	}

	resultText := attempt.EventConfig.FailText
	if passed {
		resultText = attempt.EventConfig.SuccessText
	}

	crmResult := CRMResultData{
		SessionID:   fmt.Sprintf("%d", attempt.AttemptID),
		TestID:      fmt.Sprintf("%d", attempt.EventConfig.ConfigID),
		Score:       userPoints,
		MaxScore:    maxPoints,
		IsPassed:    passed,
		CompletedAt: now.Format("2006-01-02T15:04:05Z"),
		StartedAt:   attempt.StartTime.Format("2006-01-02T15:04:05Z"),
	}

	allCompleted := checkAllTestsCompleted(claims.UserID, attempt.EventConfig)

	if allCompleted {
		if err := sendResultsToCRM(crmResult, attempt.ApplicationID); err != nil {
			fmt.Printf("Ошибка отправки результатов в CRM: %v\n", err)
		}
	}

	response := FinishAttemptResponse{
		Result:        resultText,
		Score:         userPoints,
		MaxTestPoints: maxPoints,
		Passed:        passed,
		AllCompleted:  allCompleted,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func checkAnswer(qType models.QType, options models.QuestionOptions, answer UserAnswer) bool {
	switch qType {
	case models.SingleChoice, models.MultipleChoice:
		if len(answer.Choices) != len(options.Choices) {
			return false
		}
		for i, choice := range options.Choices {
			if choice.IsTrue != answer.Choices[i] {
				return false
			}
		}
		return true

	case models.Matching:
		if len(answer.MatchingPairs) != len(options.MatchingPairs) {
			return false
		}
		for i, pair := range options.MatchingPairs {
			if pair != answer.MatchingPairs[i] {
				return false
			}
		}
		return true

	case models.CorrectOrder:
		if len(answer.Sequence) != len(options.Sequence) {
			return false
		}
		for i, item := range options.Sequence {
			if item.Order != answer.Sequence[i].Order {
				return false
			}
		}
		return true

	case models.TextInput:
		if options.CaseSensitive {
			for _, correctInput := range options.CorrectInput {
				if answer.UserInput == correctInput {
					return true
				}
			}
		} else {
			for _, correctInput := range options.CorrectInput {
				if answer.UserInput == correctInput {
					return true
				}
			}
		}
		return false
	}

	return false
}

func checkAllTestsCompleted(userID uint, currentConfig models.EventConfig) bool {
	var allMainConfigs []models.EventConfig
	database.DB.Where("event_id = ? AND specialization_id = ? AND is_extra = ?",
		currentConfig.EventID, currentConfig.SpecializationID, false).
		Find(&allMainConfigs)

	for _, cfg := range allMainConfigs {
		var attempt models.Attempt
		err := database.DB.Where("intern_id = ? AND config_id = ? AND passed = ? AND end_time IS NOT NULL",
			userID, cfg.ConfigID, true).
			First(&attempt).Error

		if err != nil {
			return false
		}
	}

	return true
}

func sendResultsToCRM(result CRMResultData, applicationID uint) error {
	cfg := config.Load()
	crmService := cfg.CRMService
	crmToken := cfg.CRMToken
	url := crmService + fmt.Sprintf("/api/users/integration/applications/%d/test-results/", applicationID)

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга данных: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(resultJSON))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", crmToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		return fmt.Errorf("CRM вернул ошибку %d: %v", resp.StatusCode, errorResponse)
	}

	fmt.Println("Результаты успешно отправлены в CRM")
	return nil
}
