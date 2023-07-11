package models

import "time"

type SensorObservation struct {
	ID                 uint    `gorm:"primary_key"`
	SensorID           uint    `gorm:"foreignKey;type:int;not null"`
	Temperature        float64 `gorm:"type:float;not null"`
	Transparency       int64   `gorm:"type:int;not null"`
	FishSpeciesID      uint    `gorm:"foreignKey;type:int;not null"`
	Count              int     `gorm:"type:int;not null"`
	DetectionTimeFrom  time.Time
	DetectionTimeUntil time.Time
}
