package models

type Sensor struct {
	ID            uint    `gorm:"primary_key"`
	Codename      string  `gorm:"type:varchar(255);not null"`
	CoordinateX   float64 `gorm:"type:float;not null"`
	CoordinateY   float64 `gorm:"type:float;not null"`
	CoordinateZ   float64 `gorm:"type:float;not null"`
	DataRate      int     `gorm:"type:int;not null"`
	SensorGroupID uint    `gorm:"foreignKey;type:int;not null"`
}
