package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"test-constructor/config"
	"test-constructor/internal/database"
	"test-constructor/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserStatisticsResponse struct {
	UserID    uint                `json:"user_id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Email     string              `json:"email"`
	Attempts  []UserAttemptDetail `json:"attempts"`
}

type UserAttemptDetail struct {
	AttemptID uint               `json:"attempt_id"`
	TestTitle string             `json:"test_title"`
	EventName string             `json:"event_name"`
	IsExtra   bool               `json:"is_extra"`
	Score     int                `json:"score"`
	MaxScore  int                `json:"max_score"`
	Passed    bool               `json:"passed"`
	Questions []QuestionStatInfo `json:"questions"`
}

type CRMEvent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetUsersResponse struct {
	Users []UserInfo `json:"users"`
}

type UserInfo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// @Summary Список стажёров
// @Security ApiKeyAuth
// @Tags manager
// @Produce json
// @Success 200 {object} GetUsersResponse
// @Failure 401 {object} map[string]string
// @Router /api/manager/users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var internRole models.Role
	if err := database.DB.Where("code = ?", "intern").First(&internRole).Error; err != nil {
		http.Error(w, "Роль стажёра не найдена", http.StatusInternalServerError)
		return
	}

	var users []models.User
	if err := database.DB.Where("role_id = ?", internRole.ID).
		Order("surname, name").
		Find(&users).Error; err != nil {
		http.Error(w, "Ошибка получения пользователей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	interns := make([]UserInfo, len(users))
	for i, user := range users {
		interns[i] = UserInfo{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
		}
	}

	response := GetUsersResponse{
		Users: interns,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Получить статистику пользователя
// @Security ApiKeyAuth
// @Description Получение статистики всех попыток конкретного пользователя
// @Tags manager
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserStatisticsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/manager/users/{id} [get]
func GetUserStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Неверный формат user_id", http.StatusBadRequest)
		return
	}

	stats, err := getUserStatistics(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка получения статистики: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func getUserStatistics(userID uint) (*UserStatisticsResponse, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var attempts []models.Attempt
	if err := database.DB.Where("intern_id = ? AND end_time IS NOT NULL", userID).
		Preload("EventConfig").
		Preload("EventConfig.Test").
		Preload("EventConfig.Test.Questions").
		Preload("Answers").
		Preload("Answers.Question").
		Order("end_time DESC").
		Find(&attempts).Error; err != nil {
		return nil, fmt.Errorf("ошибка получения попыток: %v", err)
	}

	goCFG := config.Load()
	crmService := goCFG.CRMService
	crmToken := goCFG.CRMToken

	events, err := getEventsFromCRM(crmService, crmToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения мероприятий: %v", err)
	}

	eventMap := make(map[uint]CRMEvent)
	for _, event := range events {
		eventMap[uint(event.ID)] = event
	}

	var attemptDetails []UserAttemptDetail
	for _, attempt := range attempts {
		cfg := attempt.EventConfig

		eventName := ""
		if event, exists := eventMap[cfg.EventID]; exists {
			eventName = event.Name
		} else {
			eventName = fmt.Sprintf("Event #%d", cfg.EventID)
		}

		maxScore := 0
		for _, question := range cfg.Test.Questions {
			maxScore += question.Points
		}

		questionMap := make(map[uint]models.Question)
		for _, q := range cfg.Test.Questions {
			questionMap[q.ID] = q
		}

		questions := getQuestionsStat(attempt.Answers, questionMap)

		detail := UserAttemptDetail{
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

	return &UserStatisticsResponse{
		UserID:    user.ID,
		FirstName: user.Name,
		LastName:  user.Surname,
		Email:     user.Email,
		Attempts:  attemptDetails,
	}, nil
}

func getEventsFromCRM(crmService, crmToken string) ([]CRMEvent, error) {
	req, err := http.NewRequest("GET", crmService+"/api/users/events/", nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("X-Service-Token", crmToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к CRM: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CRM вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	var events []CRMEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return events, nil
}