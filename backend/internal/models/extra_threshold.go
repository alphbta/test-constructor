package models

type ExtraThreshold struct {
	ExtraThresholdID uint `gorm:"primaryKey"`
	ConfigID         uint
	TestID           uint
	Threshold        float64 `gorm:"not null"`
	Message          string
	Config           EventConfig `gorm:"foreignKey:ConfigID;constraint:OnDelete:CASCADE;"`
}
