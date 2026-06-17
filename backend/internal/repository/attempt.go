package repository

import (
	"test-constructor/internal/domain"

	"gorm.io/gorm"
)

type AttemptRepository interface {
	Create(attempt *domain.Attempt) error
	CreateWithTx(tx *gorm.DB, attempt *domain.Attempt) error
	FindActiveByUser(userID uint) (*domain.Attempt, error)
	FindByUserAndConfig(userID, configID uint) (*domain.Attempt, error)
	Update(attempt *domain.Attempt) error
	UpdateWithTx(tx *gorm.DB, attempt *domain.Attempt) error
	FindByConfigIDsAndUser(configIDs []uint, userID uint) ([]domain.Attempt, error)
	FindByConfigID(configID uint) ([]domain.Attempt, error)
	FindCompletedByEventAndSpec(userID, eventID, specializationID uint) ([]domain.Attempt, error)
}

type attemptRepository struct {
	db *gorm.DB
}

func NewAttemptRepository(db *gorm.DB) AttemptRepository {
	return &attemptRepository{db: db}
}

func (r *attemptRepository) Create(attempt *domain.Attempt) error {
	return r.db.Create(attempt).Error
}

func (r *attemptRepository) CreateWithTx(tx *gorm.DB, attempt *domain.Attempt) error {
	return tx.Create(attempt).Error
}

func (r *attemptRepository) FindActiveByUser(userID uint) (*domain.Attempt, error) {
	var attempt domain.Attempt
	err := r.db.Where("intern_id = ? AND end_time IS NULL", userID).First(&attempt).Error
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *attemptRepository) FindByUserAndConfig(userID, configID uint) (*domain.Attempt, error) {
	var attempt domain.Attempt
	err := r.db.Where("intern_id = ? AND config_id = ?", userID, configID).First(&attempt).Error
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *attemptRepository) Update(attempt *domain.Attempt) error {
	return r.db.Save(attempt).Error
}

func (r *attemptRepository) UpdateWithTx(tx *gorm.DB, attempt *domain.Attempt) error {
	return tx.Save(attempt).Error
}

func (r *attemptRepository) FindByConfigIDsAndUser(configIDs []uint, userID uint) ([]domain.Attempt, error) {
	var attempts []domain.Attempt
	err := r.db.Where("config_id IN ? AND intern_id = ?", configIDs, userID).
		Find(&attempts).Error
	return attempts, err
}

func (r *attemptRepository) FindByConfigID(configID uint) ([]domain.Attempt, error) {
	var attempts []domain.Attempt
	err := r.db.Where("config_id = ?", configID).
		Preload("User").
		Preload("EventConfig").
		Preload("Answers").
		Preload("Answers.Question").
		Find(&attempts).Error
	return attempts, err
}

func (r *attemptRepository) FindCompletedByEventAndSpec(userID, eventID, specializationID uint) ([]domain.Attempt, error) {
	var attempts []domain.Attempt
	err := r.db.Table("attempts").
		Joins("JOIN event_configs ON event_configs.config_id = attempts.config_id").
		Where("attempts.intern_id = ? AND event_configs.event_id = ? AND event_configs.specialization_id = ? AND attempts.end_time IS NOT NULL",
			userID, eventID, specializationID).
		Find(&attempts).Error
	return attempts, err
}
