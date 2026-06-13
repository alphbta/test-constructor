package intern

import (
	"encoding/json"
	"errors"
	"net/http"
	"test-constructor/internal/auth"
	"test-constructor/internal/database"
	"test-constructor/internal/middleware"
	"test-constructor/internal/models"

	"gorm.io/gorm"
)

type UserEventCreateInfo struct {
	EventID       uint `json:"event_id"`
	ApplicationID uint `json:"application_id"`
}

type UserEventGetInfo struct {
	EventID       uint `json:"event_id"`
	UserID        uint `json:"user_id"`
	ApplicationID uint `json:"application_id"`
}

// @Summary Записаться на мероприятие
// @Security ApiKeyAuth
// @Description Создание связи между пользователем и мероприятием
// @Tags intern
// @Accept json
// @Produce json
// @Param body body UserEventCreateInfo true "Event data"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/intern/users/events [post]
func CreateUserEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var req UserEventCreateInfo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неправильный формат запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.EventID < 1 {
		http.Error(w, "Event ID должен быть положительным", http.StatusBadRequest)
		return
	}

	if req.ApplicationID < 1 {
		http.Error(w, "Application ID должен быть положительным", http.StatusBadRequest)
		return
	}

	var existingUserEvent models.UserEvent
	err := database.DB.Where("user_id = ? AND event_id = ?", claims.UserID, req.EventID).
		First(&existingUserEvent).Error

	if err == nil {
		http.Error(w, "Вы уже записаны на это мероприятие", http.StatusConflict)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userEvent := models.UserEvent{
		UserID:        claims.UserID,
		EventID:       req.EventID,
		ApplicationID: req.ApplicationID,
	}

	if err := database.DB.Create(&userEvent).Error; err != nil {
		http.Error(w, "Ошибка создания связи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := database.DB.Preload("User").First(&userEvent, userEvent.ID).Error; err != nil {
		http.Error(w, "Ошибка загрузки созданной записи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// @Summary Получить мероприятия пользователя
// @Security ApiKeyAuth
// @Description Получение списка мероприятий, на которые записался текущий пользователь
// @Tags intern
// @Produce json
// @Success 200 {array} UserEventGetInfo
// @Failure 401 {object} map[string]string
// @Router /api/intern/users/events [get]
func GetUserEvents(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var userEvents []models.UserEvent
	if err := database.DB.Preload("User").
		Where("user_id = ?", claims.UserID).
		Find(&userEvents).Error; err != nil {
		http.Error(w, "Ошибка получения данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var userEventsInfo []UserEventGetInfo
	for _, userEvent := range userEvents {
		userEventsInfo = append(userEventsInfo, UserEventGetInfo{
			EventID:       userEvent.EventID,
			UserID:        userEvent.UserID,
			ApplicationID: userEvent.ApplicationID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userEventsInfo)
}
