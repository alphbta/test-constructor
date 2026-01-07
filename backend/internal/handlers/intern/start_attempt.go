package intern

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"test-constructor/internal/auth"
	"test-constructor/internal/database"
	"test-constructor/internal/middleware"
	"test-constructor/internal/models"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type StartAttemptResponse struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Questions   []QuestionInfo `json:"questions"`
}

type QuestionInfo struct {
	QuestionID  uint          `json:"question_id"`
	Text        string        `json:"text"`
	OrderNumber int           `json:"order_number"`
	Type        models.QType  `json:"type"`
	Options     PublicOptions `json:"options"`
}

type PublicOptions struct {
	Choices       []string        `json:"choice,omitempty"`
	Matching      *PublicMatching `json:"matching,omitempty"`
	CaseSensitive bool            `json:"case_sensitive,omitempty"`
	Sequence      []string        `json:"sequence,omitempty"`
}

type PublicMatching struct {
	LeftColumn  []string `json:"left,omitempty"`
	RightColumn []string `json:"right,omitempty"`
}

// @Summary Пройти тест
// @Security ApiKeyAuth
// @Description Создание попытки
// @Tags intern
// @Accept json
// @Produce json
// @Param link path string true "Ссылка теста"
// @Success 201 {object} StartAttemptResponse
// @Router /api/intern/tests/{link} [get]
func StartAttempt(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var existingActiveAttempt models.Attempt
	if err := database.DB.Preload("Test").
		Where("intern_id = ? AND end_time IS NULL", claims.UserID).
		First(&existingActiveAttempt).Error; err == nil {
		http.Error(w, "У вас уже есть активная попытка", http.StatusConflict)
		return
	}

	vars := mux.Vars(r)
	link := vars["link"]
	var test models.Test
	err := database.DB.Preload("Questions").Where("test_link = ?", link).First(&test).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Тест не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var existingThisAttempt models.Attempt
	if err := database.DB.Preload("Test").
		Where("intern_id = ? AND test_id = ?", claims.UserID, test.ID).
		First(&existingThisAttempt).Error; err == nil {
		http.Error(w, "Вы уже прошли этот тест", http.StatusConflict)
		return
	}

	publicQuestions := make([]QuestionInfo, len(test.Questions))
	for i, q := range test.Questions {
		var options models.QuestionOptions
		if err := json.Unmarshal(q.Options, &options); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var publicOptions PublicOptions
		switch q.Type {
		case models.SingleChoice, models.MultipleChoice:
			for _, choice := range options.Choices {
				publicOptions.Choices = append(publicOptions.Choices, choice.Text)
			}
		case models.Matching:
			pairs := options.MatchingPairs
			leftColumn := make([]string, len(pairs))
			rightColumn := make([]string, len(pairs))
			rand.New(rand.NewSource(time.Now().UnixNano()))
			for i, pair := range pairs {
				leftColumn[i] = pair.LeftColumn
				rightColumn[i] = pair.RightColumn
			}

			shuffledRight := make([]string, len(rightColumn))
			copy(shuffledRight, rightColumn)
			for i := len(shuffledRight) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				shuffledRight[i], shuffledRight[j] = shuffledRight[j], shuffledRight[i]
			}

			publicOptions.Matching = &PublicMatching{leftColumn, shuffledRight}
		case models.CorrectOrder:
			rand.New(rand.NewSource(time.Now().UnixNano()))
			shuffledSequence := make([]string, len(options.Sequence))
			for i, item := range options.Sequence {
				shuffledSequence[i] = item.Text
			}

			for i := len(shuffledSequence) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				shuffledSequence[i], shuffledSequence[j] = shuffledSequence[j], shuffledSequence[i]
			}

			publicOptions.Sequence = shuffledSequence
		}

		publicQuestions[i] = QuestionInfo{
			QuestionID:  q.ID,
			Text:        q.Text,
			OrderNumber: q.OrderNumber,
			Type:        q.Type,
			Options:     publicOptions,
		}
	}

	attempt := models.Attempt{
		TestID:    test.ID,
		InternID:  claims.UserID,
		StartTime: time.Now(),
		EndTime:   nil,
	}

	if err := database.DB.Create(&attempt).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := StartAttemptResponse{
		Title:       test.Title,
		Description: test.Description,
		Questions:   publicQuestions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
