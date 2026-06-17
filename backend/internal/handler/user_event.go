package handler

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/dto"
	"test-constructor/internal/service"
)

type UserEventHandler struct {
	userEventService service.UserEventService
}

func NewUserEventHandler(userEventService service.UserEventService) *UserEventHandler {
	return &UserEventHandler{
		userEventService: userEventService,
	}
}

// CreateUserEvent записывает пользователя на мероприятие
// @Summary      Записаться на мероприятие
// @Description  Создает связь между пользователем и мероприятием
// @Security     BearerAuth
// @Tags         user-events
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUserEventRequest  true  "Данные мероприятия"
// @Success      201   {string}  string                      "Created"
// @Failure      400   {object}  dto.ErrorResponse           "Ошибка валидации"
// @Failure      409   {object}  dto.ErrorResponse           "Уже записан"
// @Router       /api/intern/users/events [post]
func (h *UserEventHandler) CreateUserEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	var req dto.CreateUserEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неправильный формат запроса")
		return
	}

	if err := h.userEventService.CreateUserEvent(claims.UserID, req); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "event ID должен быть положительным", "application ID должен быть положительным":
			status = http.StatusBadRequest
		case "вы уже записаны на это мероприятие":
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetUserEvents возвращает мероприятия пользователя
// @Summary      Получить мероприятия пользователя
// @Description  Возвращает список мероприятий, на которые записан пользователь
// @Security     BearerAuth
// @Tags         user-events
// @Produce      json
// @Success      200  {object}  dto.UserEventsListResponse  "Список мероприятий"
// @Failure      401  {object}  dto.ErrorResponse            "Не авторизован"
// @Router       /api/intern/users/events [get]
func (h *UserEventHandler) GetUserEvents(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	resp, err := h.userEventService.GetUserEvents(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка получения данных")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
