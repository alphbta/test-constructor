package handlers

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/database"
	"test-constructor/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TestResponse struct {
	ID                uint      `json:"test_id"`
	TestLink          uuid.UUID `json:"test_link"`
	CreatorID         uint      `json:"creator_id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	MarkType          int       `json:"mark_type"`
	IsProportionScore bool      `json:"is_proportion_score"`
	IsActive          bool      `json:"is_active"`
	FailText          string    `json:"fail_text"`
	SuccessText       string    `json:"success_text"`
	CompleteTime      int       `json:"complete_time"`
	Threshold         int       `json:"threshold"`
}

type TestsListResponse struct {
	Tests []TestResponse `json:"tests"`
}

func ManagerTestHandler(w http.ResponseWriter, r *http.Request) {
	var tests []models.Test
	if err := database.DB.Preload("User").Find(&tests).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response TestsListResponse
	for _, t := range tests {
		response.Tests = append(response.Tests, TestResponse{
			ID:                t.ID,
			TestLink:          t.TestLink,
			CreatorID:         t.CreatorID,
			Title:             t.Title,
			Description:       t.Description,
			MarkType:          t.MarkType,
			IsProportionScore: t.IsProportionScore,
			IsActive:          t.IsActive,
			FailText:          t.FailText,
			SuccessText:       t.SuccessText,
			CompleteTime:      t.CompleteTime,
			Threshold:         t.Threshold,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testID := vars["id"]
	if err := database.DB.Where("id = ?", testID).Delete(&models.Test{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
