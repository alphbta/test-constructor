package handlers

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/auth"
	"test-constructor/internal/database"
	"test-constructor/internal/middleware"
	"test-constructor/internal/models"
)

type AttemptInfo struct {
	AttemptID  uint   `json:"attempt_id"`
	TestTitle  string `json:"test_title"`
	ResultText string `json:"result_text"`
}

type InternAttemptResponse struct {
	AttemptsInfo []AttemptInfo `json:"attempts"`
}

func InternAttemptHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	var attempts []models.Attempt
	if err := database.DB.Preload("User").Find(&attempts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var attemptsInfo []AttemptInfo
	for _, attempt := range attempts {
		var resultText string
		if attempt.Passed {
			resultText = attempt.Test.SuccessText
		} else {
			resultText = attempt.Test.FailText
		}

		attemptInfo := AttemptInfo{
			attempt.AttemptID,
			attempt.Test.Title,
			resultText,
		}

		attemptsInfo = append(attemptsInfo, attemptInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(InternAttemptResponse{
		attemptsInfo,
	})
}
