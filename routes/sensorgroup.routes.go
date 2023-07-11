package routes

import (
	"github.com/danyaobertan/oceanstats/controllers"
	"github.com/gin-gonic/gin"
)

type SensorGroupRouteController struct {
	sensorGroupController controllers.SensorGroupController
}

func NewSensorGroupRouteController(sensorGroupController controllers.SensorGroupController) SensorGroupRouteController {
	return SensorGroupRouteController{sensorGroupController}
}

func (sgc *SensorGroupRouteController) SensorGroupRoute(rg *gin.RouterGroup) {

	router := rg.Group("/sensorgroup")
	router.POST("/", sgc.sensorGroupController.CreateSensorGroup)
	router.GET("/:id", sgc.sensorGroupController.GetSensorGroup)
	router.GET("/", sgc.sensorGroupController.GetSensorGroups)
	router.PUT("/:id", sgc.sensorGroupController.UpdateSensorGroup)
	router.DELETE("/:id", sgc.sensorGroupController.DeleteSensorGroup)
}
