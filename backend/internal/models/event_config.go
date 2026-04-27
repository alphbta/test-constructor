package models

import "github.com/google/uuid"

type EventConfig struct {
	ConfigID         uint      `gorm:"primaryKey"`
	EventID          uint      `gorm:"not null"`
	SpecializationID uint      `gorm:"not null"`
	TestID           uint      `gorm:"not null"`
	CreatorID        uint      `gorm:"not null"`
	SuccessText      string    `gorm:"not null"`
	FailureText      string    `gorm:"not null"`
	TimeLimit        int       `gorm:"default:0"`
	TestLink         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Threshold        float64   `gorm:"not null"`
	Test             Test      `gorm:"foreignKey:TestID;constraint:OnDelete:CASCADE;"`
	User             User      `gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE"`
}
