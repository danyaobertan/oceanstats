package routes

import (
	"github.com/danyaobertan/oceanstats/controllers"
	"github.com/gin-gonic/gin"
)

type SensorRouteController struct {
	sensorController controllers.SensorController
}

func NewSensorRouteController(sensorController controllers.SensorController) SensorRouteController {
	return SensorRouteController{sensorController}
}

func (sc *SensorRouteController) SensorRoute(rg *gin.RouterGroup) {

	router := rg.Group("/sensor")
	//router.POST("/", sgc.sensorController.CreateSensor)
	//router.GET("/:id", sgc.sensorController.GetSensor)
	router.GET("/", sc.sensorController.GetSensors)
	//router.PUT("/:id", sgc.sensorController.UpdateSensor)
	//router.DELETE("/:id", sgc.sensorController.DeleteSensor)
}
