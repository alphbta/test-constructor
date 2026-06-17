package repository

import (
	"test-constructor/internal/domain"

	"gorm.io/gorm"
)

type StatisticsRepository interface {
	FindInternsByRole(roleCode string) ([]domain.User, error)
	FindUserByID(userID uint) (*domain.User, error)
	FindCompletedAttemptsByUserID(userID uint) ([]domain.Attempt, error)
	FindCompletedAttemptsByConfigIDs(configIDs []uint) ([]domain.Attempt, error)
	FindConfigsByEventID(eventID uint, isExtra *bool) ([]domain.EventConfig, error)
}

type statisticsRepository struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) FindInternsByRole(roleCode string) ([]domain.User, error) {
	var role domain.Role
	if err := r.db.Where("code = ?", roleCode).First(&role).Error; err != nil {
		return nil, err
	}

	var users []domain.User
	err := r.db.Where("role_id = ?", role.ID).
		Order("surname, name").
		Find(&users).Error
	return users, err
}

func (r *statisticsRepository) FindUserByID(userID uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *statisticsRepository) FindCompletedAttemptsByUserID(userID uint) ([]domain.Attempt, error) {
	var attempts []domain.Attempt
	err := r.db.Where("intern_id = ? AND end_time IS NOT NULL", userID).
		Preload("EventConfig").
		Preload("EventConfig.Test").
		Preload("EventConfig.Test.Questions").
		Preload("Answers").
		Preload("Answers.Question").
		Order("end_time DESC").
		Find(&attempts).Error
	return attempts, err
}

func (r *statisticsRepository) FindCompletedAttemptsByConfigIDs(configIDs []uint) ([]domain.Attempt, error) {
	var attempts []domain.Attempt
	err := r.db.Where("config_id IN ? AND end_time IS NOT NULL", configIDs).
		Preload("User").
		Preload("EventConfig").
		Preload("Answers").
		Preload("Answers.Question").
		Order("end_time DESC").
		Find(&attempts).Error
	return attempts, err
}

func (r *statisticsRepository) FindConfigsByEventID(eventID uint, isExtra *bool) ([]domain.EventConfig, error) {
	query := r.db.Where("event_id = ?", eventID)

	if isExtra != nil {
		query = query.Where("is_extra = ?", *isExtra)
	}

	var configs []domain.EventConfig
	err := query.Preload("Test.Questions").Find(&configs).Error
	return configs, err
}
