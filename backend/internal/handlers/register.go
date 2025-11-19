package handlers

import (
	"encoding/json"
	"net/http"
	"test-constructor/internal/database"
	"test-constructor/internal/models"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неправильный JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" || req.Surname == "" {
		http.Error(w, "Не все поля заполнены", http.StatusBadRequest)
		return
	}

	user := models.User{
		Email:   req.Email,
		Name:    req.Name,
		Surname: req.Surname,
		Role:    models.RoleIntern,
	}

	if err := user.HashPassword(req.Password); err != nil {
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Пользователь с такой почтой уже существует", http.StatusConflict)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(RegisterResponse{
		Message: "Пользователь создан",
	})
}
