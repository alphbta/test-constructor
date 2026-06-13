package intern

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-constructor/internal/auth"
	"test-constructor/internal/database"
	"test-constructor/internal/middleware"
	"test-constructor/internal/models"
)

type TestSelectionResponse struct {
	EventID          uint       `json:"event_id"`
	SpecializationID uint       `json:"specialization_id"`
	Tests            []TestInfo `json:"tests"`
	AllCompleted     bool       `json:"all_completed"`
}

type TestInfo struct {
	ConfigID    uint   `json:"config_id"`
	TestID      uint   `json:"test_id"`
	TestLink    string `json:"test_link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeLimit   int    `json:"time_limit"`
	IsExtra     bool   `json:"is_extra"`
	Status      string `json:"status"` // "available", "locked", "in_progress", "completed"
	Score       int    `json:"score,omitempty"`
	MaxScore    int    `json:"max_score,omitempty"`
	Passed      bool   `json:"passed,omitempty"`
	AttemptID   uint   `json:"attempt_id,omitempty"`
}

// @Summary Получить список тестов для прохождения
// @Security ApiKeyAuth
// @Description Промежуточная страница со списком всех тестов в мероприятии
// @Tags intern
// @Produce json
// @Param event_id query int true "Event ID"
// @Param specialization_id query int true "Specialization ID"
// @Success 200 {object} TestSelectionResponse
// @Router /api/intern/tests/selection [get]
func GetTestSelection(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	eventIDStr := r.URL.Query().Get("event_id")

	if eventIDStr == "" {
		http.Error(w, "event_id обязателен", http.StatusBadRequest)
		return
	}

	eventID, _ := strconv.ParseUint(eventIDStr, 10, 32)

	var configs []models.EventConfig
	if err := database.DB.Preload("Test").
		Preload("ExtraThreshold").
		Where("event_id = ?", uint(eventID)).
		Find(&configs).Error; err != nil {
		http.Error(w, "Ошибка получения конфигураций", http.StatusInternalServerError)
		return
	}

	if len(configs) == 0 {
		http.Error(w, "Тесты не найдены", http.StatusNotFound)
		return
	}

	configIDs := make([]uint, len(configs))
	for i, cfg := range configs {
		configIDs[i] = cfg.ConfigID
	}

	var attempts []models.Attempt
	database.DB.Where("intern_id = ? AND config_id IN ?", claims.UserID, configIDs).
		Find(&attempts)

	attemptMap := make(map[uint]models.Attempt)
	for _, attempt := range attempts {
		attemptMap[attempt.ConfigID] = attempt
	}

	mainConfigs := make(map[uint]bool)
	for _, cfg := range configs {
		if !cfg.IsExtra {
			mainConfigs[cfg.ConfigID] = true
		}
	}

	var tests []TestInfo
	allCompleted := true

	for _, config := range configs {
		attempt, hasAttempt := attemptMap[config.ConfigID]

		status := "available"

		if config.IsExtra {
			status = "locked"
			if hasAccessToExtraConfig(claims.UserID, config.ConfigID) {
				status = "available"
			}
		}

		testInfo := TestInfo{
			ConfigID:    config.ConfigID,
			TestID:      config.TestID,
			TestLink:    config.TestLink.String(),
			Title:       config.Test.Title,
			Description: config.Test.Description,
			TimeLimit:   config.TimeLimit,
			IsExtra:     config.IsExtra,
			Status:      status,
		}

		if hasAttempt {
			if attempt.EndTime != nil {
				testInfo.Status = "completed"
				testInfo.Score = attempt.Score
				testInfo.MaxScore = attempt.MaxScore
				testInfo.Passed = attempt.Passed
				testInfo.AttemptID = attempt.AttemptID
			} else {
				testInfo.Status = "in_progress"
				testInfo.AttemptID = attempt.AttemptID
				allCompleted = false
			}
		} else {
			allCompleted = false
		}

		tests = append(tests, testInfo)
	}

	response := TestSelectionResponse{
		EventID:      uint(eventID),
		Tests:        tests,
		AllCompleted: allCompleted,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func hasAccessToExtraConfig(userID uint, configID uint) bool {
	var extraThreshold models.ExtraThreshold
	if err := database.DB.Where("extra_config_id = ?", configID).
		First(&extraThreshold).Error; err != nil {
		return false
	}

	var attempt models.Attempt
	if err := database.DB.Where("intern_id = ? AND config_id = ? AND end_time IS NOT NULL",
		userID, extraThreshold.ConfigID).
		First(&attempt).Error; err != nil {
		return false
	}

	return attempt.Score >= extraThreshold.Threshold
}
