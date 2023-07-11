package models

type SensorGroup struct {
	ID   uint   `gorm:"primary_key" json:"id,omitempty"`
	Name string `gorm:"type:varchar(255);uniqueIndex;not null" json:"title,omitempty"`
}

type CreateSensorGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateSensorGroupRequest struct {
	Name string `json:"name,omitempty"`
}
