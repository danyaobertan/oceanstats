package controllers

import (
	"github.com/danyaobertan/oceanstats/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type SensorGroupController struct {
	DB *gorm.DB
}

func NewSensorGroupController(DB *gorm.DB) SensorGroupController {
	return SensorGroupController{DB}
}

func (sgc *SensorGroupController) CreateSensorGroup(ctx *gin.Context) {
	var payload *models.CreateSensorGroupRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	NewSensorGroup := models.SensorGroup{
		Name: payload.Name,
	}

	result := sgc.DB.Create(&NewSensorGroup)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "SensorGroup with that name already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": NewSensorGroup})
}

func (sgc *SensorGroupController) UpdateSensorGroup(ctx *gin.Context) {
	sensorGroupId := ctx.Param("id")

	var payload *models.UpdateSensorGroupRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var updatedSensorGroup models.SensorGroup
	result := sgc.DB.First(&updatedSensorGroup, "id = ?", sensorGroupId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No SensorGroup with that name exists"})
		return
	}

	sensorGroupToUpdate := models.SensorGroup{
		Name: payload.Name,
	}

	sgc.DB.Model(&updatedSensorGroup).Updates(sensorGroupToUpdate)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedSensorGroup})
}

func (sgc *SensorGroupController) GetSensorGroup(ctx *gin.Context) {
	sensorGroupId := ctx.Param("id")

	var sensorGroup models.SensorGroup
	result := sgc.DB.First(&sensorGroup, "id = ?", sensorGroupId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No SensorGroup with that name exists"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": sensorGroup})
}

func (sgc *SensorGroupController) GetSensorGroups(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var sensorGroups []models.SensorGroup
	result := sgc.DB.Limit(intLimit).Offset(offset).Find(&sensorGroups)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No SensorGroups found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": sensorGroups})
}

func (sgc *SensorGroupController) DeleteSensorGroup(ctx *gin.Context) {
	sensorGroupId := ctx.Param("id")

	result := sgc.DB.Delete(&models.SensorGroup{}, "id = ?", sensorGroupId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No SensorGroup with that name exists"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "SensorGroup deleted successfully"})
}
