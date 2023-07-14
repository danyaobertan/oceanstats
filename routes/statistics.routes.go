package routes

import (
	"github.com/danyaobertan/oceanstats/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StatisticsRouteController struct {
	statisticsController controllers.StatisticsController
}

func NewStatisticsRouteController(statisticsController controllers.StatisticsController) StatisticsRouteController {
	return StatisticsRouteController{statisticsController}
}

func (sc *StatisticsRouteController) StatisticsRoute(rg *gin.RouterGroup) {
	router := rg.Group("/statistics")
	router.GET("/group/:groupName/transparency/average", sc.statisticsController.GetGroupAverageTransparency)
	router.GET("/group/:groupName/temperature/average", sc.statisticsController.GetGroupAverageTemperature)

	router.GET("/group/:groupName/species", sc.statisticsController.GetGroupSpecies)

	//router.GET("/group/:groupName/species/top/:N", sc.statisticsController.GetTopGroupSpecies)

	router.GET("/group/:groupName/species/top/:N", func(ctx *gin.Context) {
		groupName := ctx.Param("groupName")
		N := ctx.Param("N")
		from := ctx.Query("from")
		till := ctx.Query("till")
		if from != "" && till != "" {
			topSpecies := sc.statisticsController.GetTopGroupSpeciesBetweenDates(groupName, N, from, till)
			ctx.JSON(http.StatusOK, gin.H{
				"top_species_between_dates": topSpecies,
			})
			return
		} else {
			topSpecies := sc.statisticsController.GetTopGroupSpecies(groupName, N)
			ctx.JSON(http.StatusOK, gin.H{
				"top_species": topSpecies,
			})
		}

	})

	router.GET("/region/temperature/min", func(ctx *gin.Context) {
		xMin := ctx.Query("xMin")
		xMax := ctx.Query("xMax")
		yMin := ctx.Query("yMin")
		yMax := ctx.Query("yMax")
		zMin := ctx.Query("zMin")
		zMax := ctx.Query("zMax")
		minTemperature := sc.statisticsController.GetRegionMinTemperature(xMin, xMax, yMin, yMax, zMin, zMax)
		ctx.JSON(http.StatusOK, gin.H{
			"minimum_temperature": minTemperature,
		})
	})

	router.GET("/region/temperature/max", func(ctx *gin.Context) {
		xMin := ctx.Query("xMin")
		xMax := ctx.Query("xMax")
		yMin := ctx.Query("yMin")
		yMax := ctx.Query("yMax")
		zMin := ctx.Query("zMin")
		zMax := ctx.Query("zMax")
		maxTemperature := sc.statisticsController.GetRegionMaxTemperature(xMin, xMax, yMin, yMax, zMin, zMax)
		ctx.JSON(http.StatusOK, gin.H{
			"maximum_temperature": maxTemperature,
		})
	})

	router.GET("/sensor/:codeName/temperature/average", func(ctx *gin.Context) {
		codeName := ctx.Param("codeName")
		from := ctx.Query("from")
		till := ctx.Query("till")
		maxTemperature := sc.statisticsController.GetSensorTemperatureAverageBetweenDates(codeName, from, till)
		if maxTemperature <= -273.0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no data for the specified time range",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"average_temperature_between_dates": maxTemperature,
		})
	})
}
