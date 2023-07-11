package models

type FishSpecies struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"type:varchar(255);not null"`
}
