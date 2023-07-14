package routes

import (
	"github.com/danyaobertan/oceanstats/controllers"
	"github.com/gin-gonic/gin"
)

type SensorObservationRouteController struct {
	sensorObservationController controllers.SensorObservationController
}

func NewSensorObservationRouteController(sensorObservationController controllers.SensorObservationController) SensorObservationRouteController {
	return SensorObservationRouteController{sensorObservationController}
}

func (soc *SensorObservationRouteController) SensorObservationRoute(rg *gin.RouterGroup) {

	router := rg.Group("/observation")
	router.POST("/", soc.sensorObservationController.CreateSensorObservation)
	router.POST("/bulk", soc.sensorObservationController.CreateBulkSensorObservationHandler)

	//router.GET("/:id", soc.sensorObservationController.GetSensorObservation)
	//router.GET("/", soc.sensorObservationController.GetSensorObservations)
	//router.PUT("/:id", soc.sensorObservationController.UpdateSensorObservation)
	//router.DELETE("/:id", soc.sensorObservationController.DeleteSensorObservation)
}
