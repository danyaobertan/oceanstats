package models

import "time"

type SensorObservation struct {
	ID                 uint    `gorm:"primary_key" json:"id,omitempty"`
	SensorID           uint    `gorm:"foreignKey;type:int;not null"`
	Temperature        float64 `gorm:"type:float;not null"`
	Transparency       int64   `gorm:"type:int;not null"`
	FishSpeciesID      uint    `gorm:"foreignKey;type:int;not null"`
	Count              int     `gorm:"type:int;not null"`
	DetectionTimeFrom  time.Time
	DetectionTimeUntil time.Time
}

type CreateSensorObservationRequest struct {
	SensorID           uint      `json:"sensor_id" binding:"required"`
	Temperature        float64   `json:"temperature" binding:"required"`
	Transparency       int64     `json:"transparency" binding:"required"`
	FishSpeciesID      uint      `json:"fish_species_id" binding:"required"`
	Count              int       `json:"count" binding:"required"`
	DetectionTimeFrom  time.Time `json:"detection_time_from" binding:"required"`
	DetectionTimeUntil time.Time `json:"detection_time_until" binding:"required"`
}
