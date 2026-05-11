package models

type ExtraThreshold struct {
	ExtraThresholdID uint `gorm:"primaryKey"`
	ConfigID         uint `gorm:"not null"`
	TestID           uint
	Threshold        float64 `gorm:"not null"`
	Message          string
}
