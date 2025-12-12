package manager

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/database"
	"test-constructor/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TestInfo struct {
	ID           uint      `json:"test_id"`
	TestLink     uuid.UUID `json:"test_link"`
	CreatorID    uint      `json:"creator_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	IsPercentage bool      `json:"is_percentage"`
	IsActive     bool      `json:"is_active"`
	FailText     string    `json:"fail_text"`
	SuccessText  string    `json:"success_text"`
	CompleteTime int       `json:"complete_time"`
	Threshold    int       `json:"threshold"`
}

type TestsInfoResponse struct {
	Tests []TestInfo `json:"tests"`
}

// @Summary Получить тесты для организотора
// @Security ApiKeyAuth
// @Description Получить список тестов
// @Tags manager
// @Produce json
// @Success 200 {object} TestsInfoResponse
// @Failure 404 {object} map[string]string
// @Router /api/manager/tests [get]
func ManagerTestHandler(w http.ResponseWriter, r *http.Request) {
	var tests []models.Test
	if err := database.DB.Preload("User").Find(&tests).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response TestsInfoResponse
	for _, t := range tests {
		response.Tests = append(response.Tests, TestInfo{
			ID:           t.ID,
			TestLink:     t.TestLink,
			CreatorID:    t.CreatorID,
			Title:        t.Title,
			Description:  t.Description,
			IsPercentage: t.IsPercentage,
			IsActive:     t.IsActive,
			FailText:     t.FailText,
			SuccessText:  t.SuccessText,
			CompleteTime: t.CompleteTime,
			Threshold:    t.Threshold,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Удалить тест
// @Security ApiKeyAuth
// @Tags manager
// @Accept json
// @Produce json
// @Param id path int true "ID теста" minimum(1)
// @Success 200
// @Router /api/manager/tests/delete/{id} [post]
func DeleteTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testID := vars["id"]
	if err := database.DB.Where("id = ?", testID).Delete(&models.Test{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
