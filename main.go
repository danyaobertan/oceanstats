package main

import (
	"github.com/danyaobertan/oceanstats/controllers"
	_ "github.com/danyaobertan/oceanstats/docs"
	"github.com/danyaobertan/oceanstats/initializers"
	"github.com/danyaobertan/oceanstats/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

var (
	server *gin.Engine

	SensorGroupController       controllers.SensorGroupController
	SensorGroupRouterController routes.SensorGroupRouteController

	SensorController       controllers.SensorController
	SensorRouterController routes.SensorRouteController

	SensorObservationController       controllers.SensorObservationController
	SensorObservationRouterController routes.SensorObservationRouteController

	StatisticsController       controllers.StatisticsController
	StatisticsRouterController routes.StatisticsRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectPostgres(&config)
	initializers.ConnectRedis(&config)

	SensorGroupController = controllers.NewSensorGroupController(initializers.DB)
	SensorGroupRouterController = routes.NewSensorGroupRouteController(SensorGroupController)

	SensorController = controllers.NewSensorController(initializers.DB)
	SensorRouterController = routes.NewSensorRouteController(SensorController)

	SensorObservationController = controllers.NewSensorObservationController(initializers.DB)
	SensorObservationRouterController = routes.NewSensorObservationRouteController(SensorObservationController)

	StatisticsController = controllers.NewStatisticsController(initializers.DB, initializers.RedisClient)
	StatisticsRouterController = routes.NewStatisticsRouteController(StatisticsController)

	server = gin.Default()
}

// @title           OceanStats API (Go language test task)
// @version         1.0
// @description     OceanStats is a backend API implementation for a set of underwater sensors developed as a test task for Helo Labs.
// @termsOfService  https://github.com/danyaobertan/oceanstats

// @contact.name   Danyil-Mykola Obertan
// @contact.email  danyilmykolaobertan@gmail.com

// @host      localhost:8000
// @BasePath  /api/
func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Golang language test task"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	SensorRouterController.SensorRoute(router)
	SensorGroupRouterController.SensorGroupRoute(router)
	SensorObservationRouterController.SensorObservationRoute(router)
	StatisticsRouterController.StatisticsRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
