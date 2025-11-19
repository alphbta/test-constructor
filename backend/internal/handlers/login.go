package handlers

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/database"
	"test-constructor/internal/models"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неправильный JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Не все поля заполнены", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Вы вошли",
	})
}
