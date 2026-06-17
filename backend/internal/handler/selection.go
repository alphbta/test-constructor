package handler

import (
	"net/http"
	"strconv"
	"test-constructor/internal/service"
)

type TestSelectionHandler struct {
	testSelectionService service.TestSelectionService
}

func NewTestSelectionHandler(
	testSelectionService service.TestSelectionService,
) *TestSelectionHandler {
	return &TestSelectionHandler{
		testSelectionService: testSelectionService,
	}
}

// GetTestSelection возвращает список тестов мероприятия
// @Summary      Получить список тестов
// @Description  Возвращает список доступных тестов для мероприятия с учётом замен
// @Security     BearerAuth
// @Tags         attempts
// @Produce      json
// @Param        event_id  query     int                          true  "Event ID"
// @Success      200       {object}  dto.TestSelectionResponse    "Список тестов"
// @Failure      400       {object}  dto.ErrorResponse            "Ошибка валидации"
// @Failure      401       {object}  dto.ErrorResponse            "Не авторизован"
// @Failure      404       {object}  dto.ErrorResponse            "Пользователь не записан на мероприятие"
// @Router       /api/intern/tests/selection [get]
func (h *TestSelectionHandler) GetTestSelection(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	eventIDStr := r.URL.Query().Get("event_id")
	if eventIDStr == "" {
		writeError(w, http.StatusBadRequest, "event_id обязателен")
		return
	}

	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат event_id")
		return
	}

	resp, err := h.testSelectionService.GetSelection(claims.UserID, uint(eventID))
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "вы не записаны на это мероприятие":
			status = http.StatusNotFound
		default:
			status = http.StatusBadRequest
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
