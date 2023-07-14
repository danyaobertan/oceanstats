package controllers

import (
	"github.com/danyaobertan/oceanstats/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type SensorController struct {
	DB *gorm.DB
}

func NewSensorController(DB *gorm.DB) SensorController {
	return SensorController{DB}
}

func (sc *SensorController) GetSensors(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "1000")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var sensors []models.Sensor
	result := sc.DB.Limit(intLimit).Offset(offset).Find(&sensors)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No SensorGroups found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": sensors})
}
