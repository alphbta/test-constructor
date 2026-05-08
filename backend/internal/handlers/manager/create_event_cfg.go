package manager

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/auth"
	"test-constructor/internal/database"
	"test-constructor/internal/middleware"
	"test-constructor/internal/models"
)

type CreateEventCfgInfo struct {
	EventID          uint                 `json:"event_id"`
	SpecializationID uint                 `json:"specialization_id"`
	TestID           uint                 `json:"test_id"`
	SuccessText      string               `json:"success_text"`
	FailText         string               `json:"fail_text"`
	TimeLimit        int                  `json:"time_limit"`
	Threshold        float64              `json:"threshold"`
	ExtraThreshold   []ExtraThresholdInfo `json:"extra_threshold"`
}

type ExtraThresholdInfo struct {
	Threshold float64 `json:"threshold"`
	Message   string  `json:"message"`
	TestID    uint    `json:"test_id"`
}

func CreateConfig(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok {
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}

	var req CreateEventCfgInfo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неправильный JSON", http.StatusBadRequest)
		return
	}

	if req.EventID < 1 || req.SpecializationID < 1 || req.TestID < 1 {
		http.Error(w, "ID должен быть положительным", http.StatusBadRequest)
		return
	}

	if req.Threshold < 1 {
		http.Error(w, "Пороговое значение должно быть положительным", http.StatusBadRequest)
	}

	userID := claims.UserID
	transaction := database.DB.Begin()
	if transaction.Error != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
	}()

	eventCFG := models.EventConfig{
		EventID:          req.EventID,
		SpecializationID: req.SpecializationID,
		TestID:           req.TestID,
		CreatorID:        userID,
		SuccessText:      req.SuccessText,
		FailText:         req.FailText,
		TimeLimit:        req.TimeLimit,
		Threshold:        req.Threshold,
	}

	if err := transaction.Create(&eventCFG).Error; err != nil {
		transaction.Rollback()
		http.Error(w, "Ошибка создания настройки: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, eThreshold := range req.ExtraThreshold {
		extraThreshold := models.ExtraThreshold{
			Threshold: eThreshold.Threshold,
			Message:   eThreshold.Message,
			TestID:    eThreshold.TestID,
		}

		if err := transaction.Create(&extraThreshold).Error; err != nil {
			transaction.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}
