package controllers

import (
	"github.com/danyaobertan/oceanstats/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type SensorObservationController struct {
	DB *gorm.DB
}

func NewSensorObservationController(DB *gorm.DB) SensorObservationController {
	return SensorObservationController{DB}
}

func (soc *SensorObservationController) CreateSensorObservation(ctx *gin.Context) {
}

func (soc *SensorObservationController) CreateBulkSensorObservationHandler(ctx *gin.Context) {

	var request struct {
		Observations []models.CreateSensorObservationRequest `json:"observations"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := soc.CreateBulkSensorObservation(request.Observations)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bulk insertion successful"})

}

func (soc *SensorObservationController) CreateBulkSensorObservation(observations []models.CreateSensorObservationRequest) error {
	batchSize := 500 // Number of records to insert in each batch

	// Calculate the number of batches
	numBatches := len(observations) / batchSize
	if len(observations)%batchSize != 0 {
		numBatches++
	}

	// Insert records in batches
	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := start + batchSize

		if end > len(observations) {
			end = len(observations)
		}

		batch := observations[start:end]

		// Convert CreateSensorObservationRequest to SensorObservation
		var sensorObservations []models.SensorObservation
		for _, obs := range batch {
			sensorObservation := models.SensorObservation{
				SensorID:           obs.SensorID,
				Temperature:        obs.Temperature,
				Transparency:       obs.Transparency,
				FishSpeciesID:      obs.FishSpeciesID,
				Count:              obs.Count,
				DetectionTimeFrom:  obs.DetectionTimeFrom,
				DetectionTimeUntil: obs.DetectionTimeUntil,
			}
			sensorObservations = append(sensorObservations, sensorObservation)
		}

		// Perform bulk insertion
		if err := soc.DB.CreateInBatches(sensorObservations, len(sensorObservations)).Error; err != nil {
			return err
		}
	}

	return nil
}
