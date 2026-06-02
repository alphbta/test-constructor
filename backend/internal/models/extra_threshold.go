package models

type ExtraThreshold struct {
	ExtraThresholdID uint `gorm:"primaryKey"`
	ConfigID         uint `gorm:"not null"`
	ExtraConfigID    uint `gorm:"not null"`
	Threshold        int  `gorm:"not null"`
	Message          string
	ExtraConfig      EventConfig `gorm:"foreignKey:ExtraConfigID;constraint:OnDelete:CASCADE;"`
}
