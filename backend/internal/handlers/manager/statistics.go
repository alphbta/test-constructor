package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"test-constructor/internal/database"
	"test-constructor/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type StatisticsRequest struct {
	IsExtra *bool `json:"is_extra,omitempty"`
}

type StatisticsResponse struct {
	Attempts []UserAttemptInfo `json:"attempts"`
}

type UserAttemptInfo struct {
	UserID    uint               `json:"user_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Score     int                `json:"score"`
	MaxScore  int                `json:"max_score"`
	Passed    bool               `json:"passed"`
	TimeSpent int                `json:"time_spent_minutes"`
	IsExtra   bool               `json:"is_extra"`
	Questions []QuestionStatInfo `json:"questions"`
}

type QuestionStatInfo struct {
	Text         string  `json:"text"`
	Points       float64 `json:"points_earned"`
	MaxPoints    int     `json:"max_points"`
	IsCorrect    bool    `json:"is_correct"`
	QuestionType string  `json:"question_type"`
	OrderNumber  int     `json:"order_number"`
}

// @Summary Получить статистику по мероприятию
// @Security ApiKeyAuth
// @Description Получение статистики всех попыток по мероприятию
// @Tags manager
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param is_extra body bool false "Фильтр по дополнительным тестам"
// @Success 200 {object} StatisticsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/manager/events/{id}/statistics [get]
func GetEventStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventIDStr := vars["id"]

	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Неверный формат event_id", http.StatusBadRequest)
		return
	}

	var filter StatisticsRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			filter = StatisticsRequest{}
		}
	}

	stats, err := getEventStatistics(uint(eventID), filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Мероприятие не найдено", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка получения статистики: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if len(stats.Attempts) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StatisticsResponse{
			Attempts: []UserAttemptInfo{},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func getEventStatistics(eventID uint, filter StatisticsRequest) (*StatisticsResponse, error) {
	var configs []models.EventConfig
	query := database.DB.Where("event_id = ?", eventID)

	if filter.IsExtra != nil {
		query = query.Where("is_extra = ?", *filter.IsExtra)
	}

	if err := query.Preload("Test.Questions").Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("ошибка получения конфигураций: %v", err)
	}

	if len(configs) == 0 {
		return &StatisticsResponse{
			Attempts: []UserAttemptInfo{},
		}, nil
	}

	configIDs := make([]uint, len(configs))
	configMap := make(map[uint]models.EventConfig)
	questionMap := make(map[uint]map[uint]models.Question)

	for i, cfg := range configs {
		configIDs[i] = cfg.ConfigID
		configMap[cfg.ConfigID] = cfg

		questionMap[cfg.ConfigID] = make(map[uint]models.Question)
		for _, q := range cfg.Test.Questions {
			questionMap[cfg.ConfigID][q.ID] = q
		}
	}

	var attempts []models.Attempt
	if err := database.DB.Where("config_id IN ? AND end_time IS NOT NULL", configIDs).
		Preload("User").
		Preload("EventConfig").
		Preload("Answers").
		Preload("Answers.Question").
		Order("end_time DESC").
		Find(&attempts).Error; err != nil {
		return nil, fmt.Errorf("ошибка получения попыток: %v", err)
	}

	var userAttempts []UserAttemptInfo
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

		questions := getQuestionsStat(attempt.Answers, questionMap[cfg.ConfigID])

		attemptInfo := UserAttemptInfo{
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

	return &StatisticsResponse{
		Attempts: userAttempts,
	}, nil
}

func getQuestionsStat(answers []models.Answer, questionsMap map[uint]models.Question) []QuestionStatInfo {
	var questionsStat []QuestionStatInfo

	for _, answer := range answers {
		question, exists := questionsMap[answer.QuestionID]
		if !exists {
			if answer.Question.ID != 0 {
				question = answer.Question
			} else {
				continue
			}
		}

		questionInfo := QuestionStatInfo{
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