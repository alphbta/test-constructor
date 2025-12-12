package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	TestLink     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"test_link"`
	CreatorID    uint      `json:"creator_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	IsPercentage bool      `json:"is_percentage"`
	IsActive     bool      `json:"is_active"`
	FailText     string    `json:"fail_text"`
	SuccessText  string    `json:"success_text"`
	CompleteTime int       `json:"complete_time"`
	Threshold    int       `json:"threshold"`
	// Связи
	Questions []Question
	User      User `gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE"`
}
