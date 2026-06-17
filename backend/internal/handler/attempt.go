package handler

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/dto"
	"test-constructor/internal/service"

	"github.com/gorilla/mux"
)

type AttemptHandler struct {
	attemptService service.AttemptService
}

func NewAttemptHandler(
	attemptService service.AttemptService,
) *AttemptHandler {
	return &AttemptHandler{
		attemptService: attemptService,
	}
}

// StartAttempt начинает попытку прохождения теста
// @Summary      Начать тест
// @Description  Создаёт новую попытку прохождения теста по ссылке конфигурации
// @Security     BearerAuth
// @Tags         attempts
// @Accept       json
// @Produce      json
// @Param        link  path      string                      true  "Ссылка конфигурации теста (UUID)"
// @Param        body  body      dto.StartAttemptRequest     true  "Данные для начала попытки"
// @Success      201   {object}  dto.StartAttemptResponse    "Попытка создана"
// @Failure      400   {object}  dto.ErrorResponse           "Ошибка валидации"
// @Failure      401   {object}  dto.ErrorResponse           "Не авторизован"
// @Failure      403   {object}  dto.ErrorResponse           "Нет доступа к тесту"
// @Failure      404   {object}  dto.ErrorResponse           "Тест не найден"
// @Failure      409   {object}  dto.ErrorResponse           "Уже есть активная попытка или тест пройден"
// @Router       /api/intern/tests/{link} [get]
func (h *AttemptHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	vars := mux.Vars(r)
	link := vars["link"]
	if link == "" {
		writeError(w, http.StatusBadRequest, "Не указана ссылка на тест")
		return
	}

	var req dto.StartAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неправильный формат запроса")
		return
	}

	if req.ApplicationID == 0 {
		writeError(w, http.StatusBadRequest, "Application ID обязателен")
		return
	}

	resp, err := h.attemptService.StartAttempt(claims.UserID, link, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "тест не найден":
			status = http.StatusNotFound
		case "у вас уже есть активная попытка. Завершите её перед началом новой":
			status = http.StatusConflict
		case "вы уже прошли этот тест":
			status = http.StatusConflict
		case "у вас нет доступа к этому дополнительному тесту":
			status = http.StatusForbidden
		default:
			status = http.StatusBadRequest
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetActiveAttempt возвращает информацию об активной попытке
// @Summary      Получить активную попытку
// @Description  Возвращает данные текущей незавершённой попытки пользователя
// @Security     BearerAuth
// @Tags         attempts
// @Produce      json
// @Success      200  {object}  dto.StartAttemptResponse  "Активная попытка"
// @Failure      401  {object}  dto.ErrorResponse          "Не авторизован"
// @Failure      404  {object}  dto.ErrorResponse          "Нет активной попытки"
// @Router       /api/intern/attempt/active [get]
func (h *AttemptHandler) GetActiveAttempt(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	resp, err := h.attemptService.GetActiveAttempt(claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "активная попытка не найдена" {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// FinishAttempt завершает попытку и проверяет ответы
// @Summary      Завершить тест
// @Description  Проверяет ответы пользователя и завершает активную попытку
// @Security     BearerAuth
// @Tags         attempts
// @Accept       json
// @Produce      json
// @Param        answers  body      dto.FinishAttemptRequest    true  "Ответы пользователя"
// @Success      200      {object}  dto.FinishAttemptResponse   "Результаты проверки"
// @Failure      400      {object}  dto.ErrorResponse           "Ошибка валидации"
// @Failure      401      {object}  dto.ErrorResponse           "Не авторизован"
// @Failure      404      {object}  dto.ErrorResponse           "Активная попытка не найдена"
// @Router       /api/intern/attempt/finish [post]
func (h *AttemptHandler) FinishAttempt(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	var req dto.FinishAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неправильный формат запроса")
		return
	}

	if len(req.UserAnswers) == 0 {
		writeError(w, http.StatusBadRequest, "Не переданы ответы")
		return
	}

	resp, err := h.attemptService.FinishAttempt(claims.UserID, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "активная попытка не найдена":
			status = http.StatusNotFound
		case "ответы не соответствуют тесту":
			status = http.StatusBadRequest
		default:
			if len(err.Error()) > 0 {
				status = http.StatusBadRequest
			}
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
