package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-constructor/internal/dto"
	"test-constructor/internal/service"

	"github.com/gorilla/mux"
)

type EventConfigHandler struct {
	eventConfigService service.EventConfigService
}

func NewEventConfigHandler(eventConfigService service.EventConfigService) *EventConfigHandler {
	return &EventConfigHandler{
		eventConfigService: eventConfigService,
	}
}

// CreateConfig создает конфигурацию мероприятия
// @Summary      Создать конфигурацию мероприятия
// @Description  Создает новую конфигурацию тестирования для мероприятия
// @Security     BearerAuth
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        config  body      dto.CreateEventConfigRequest    true  "Конфигурация"
// @Param specialization_id body uint false "Specialization ID (0 = общий тест)"
// @Success      201     {object}  dto.CreateEventConfigResponse   "Конфигурация создана"
// @Failure      400     {object}  dto.ErrorResponse               "Ошибка валидации"
// @Failure      401     {object}  dto.ErrorResponse               "Не авторизован"
// @Failure      500     {object}  dto.ErrorResponse               "Внутренняя ошибка"
// @Router       /api/manager/events [post]
func (h *EventConfigHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Вы не авторизованы")
		return
	}

	var req dto.CreateEventConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неправильный JSON")
		return
	}

	resp, err := h.eventConfigService.CreateConfig(claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// UpdateConfig обновляет конфигурацию мероприятия
// @Summary      Обновить конфигурацию
// @Description  Обновляет существующую конфигурацию мероприятия
// @Security     BearerAuth
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id      path      int                             true  "Config ID"
// @Param        config  body      dto.UpdateEventConfigRequest    true  "Обновленная конфигурация"
// @Success      200     {object}  dto.UpdateEventConfigResponse   "Конфигурация обновлена"
// @Failure      400     {object}  dto.ErrorResponse               "Ошибка валидации"
// @Failure      403     {object}  dto.ErrorResponse               "Нет прав"
// @Failure      404     {object}  dto.ErrorResponse               "Не найдена"
// @Router       /api/manager/events/{id} [put]
func (h *EventConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Вы не авторизованы")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный ID")
		return
	}

	var req dto.UpdateEventConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неправильный JSON")
		return
	}

	resp, err := h.eventConfigService.UpdateConfig(uint(id), claims.UserID, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "конфигурация не найдена":
			status = http.StatusNotFound
		case "у вас нет прав на редактирование":
			status = http.StatusForbidden
		default:
			status = http.StatusBadRequest
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
