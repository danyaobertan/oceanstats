package main

import (
	"fmt"
	"github.com/danyaobertan/oceanstats/initializers"
	"github.com/danyaobertan/oceanstats/models"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	err := initializers.DB.AutoMigrate(&models.SensorGroup{}, &models.Sensor{}, &models.FishSpecies{}, &models.SensorObservation{})
	if err != nil {
		log.Fatal("Could not migrate the Database", err)
		return
	}
	fmt.Println("Migration complete")
}
