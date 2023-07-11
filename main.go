package main

import (
	"github.com/danyaobertan/oceanstats/controllers"
	"github.com/danyaobertan/oceanstats/initializers"
	"github.com/danyaobertan/oceanstats/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	server *gin.Engine

	SensorGroupController       controllers.SensorGroupController
	SensorGroupRouterController routes.SensorGroupRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	SensorGroupController = controllers.NewSensorGroupController(initializers.DB)
	SensorGroupRouterController = routes.NewSensorGroupRouteController(SensorGroupController)

	server = gin.Default()
}

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

	SensorGroupRouterController.SensorGroupRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
