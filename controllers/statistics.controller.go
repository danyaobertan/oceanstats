package controllers

import (
	"database/sql"
	"fmt"
	"github.com/danyaobertan/oceanstats/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type StatisticsController struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func NewStatisticsController(DB *gorm.DB, RedisClient *redis.Client) StatisticsController {
	return StatisticsController{DB, RedisClient}
}

type SpeciesCount struct {
	Name  string `gorm:"column:name"`
	Count int    `gorm:"column:species_count"`
}

// GetGroupAverageTransparency  		   godoc
// @Summary      Get current average transparency inside the group
// @Description  Responds with the current average transparency inside the group
// @Tags         statistics
// @Produce      json
// @Router	    /statistics/group/{groupName}/transparency/average [get]
// @Param groupName path string true "Name of the group" example(alpha) default(alpha)
func (sc *StatisticsController) GetGroupAverageTransparency(ctx *gin.Context) {

	groupName := ctx.Param("groupName")

	var result struct {
		AverageTransparency float64 `gorm:"column:average_transparency"`
	}
	var sensorObservation models.SensorObservation

	err := sc.DB.Model(&sensorObservation).
		Select("AVG(sensor_observations.transparency) AS average_transparency").
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Joins("JOIN sensor_groups ON sensors.sensor_group_id = sensor_groups.id").
		Where("sensor_groups.name = ?", groupName).
		Scan(&result).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve average transparency"})
		return
	}
	cacheKey := "transparency_average:" + groupName
	cacheValue := fmt.Sprintf("%.2f", result.AverageTransparency)

	err = sc.RedisClient.Set(ctx, cacheKey, cacheValue, 10*time.Second).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key | Error: %v - cacheKey: %s - cacheValue: %s\n", err, cacheKey, cacheValue))
	}
	ctx.JSON(http.StatusOK, gin.H{"average_transparency": result.AverageTransparency})
}

// GetGroupAverageTemperature godoc
// @Summary      Get current average temperature inside the group
// @Description  Responds with the current average temperature inside the group
// @Tags         statistics
// @Produce      json
// @Router	    /statistics/group/{groupName}/temperature/average [get]
// @Param groupName path string true "Name of the group" example(alpha) default(alpha)
func (sc *StatisticsController) GetGroupAverageTemperature(ctx *gin.Context) {
	groupName := ctx.Param("groupName")

	var result struct {
		AverageTemperature float64 `gorm:"column:average_temperature"`
	}
	var sensorObservation models.SensorObservation

	err := sc.DB.Model(&sensorObservation).
		Select("AVG(sensor_observations.temperature) AS average_temperature").
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Joins("JOIN sensor_groups ON sensors.sensor_group_id = sensor_groups.id").
		Where("sensor_groups.name = ?", groupName).
		Scan(&result).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve average temperature"})
		return
	}
	cacheKey := "temperature_average:" + groupName
	cacheValue := fmt.Sprintf("%.2f", result.AverageTemperature)
	err = sc.RedisClient.Set(ctx, cacheKey, cacheValue, 10*time.Second).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key | Error: %v - cacheKey: %s - cacheValue: %s\n", err, cacheKey, cacheValue))
	}
	ctx.JSON(http.StatusOK, gin.H{"average_temperature": result.AverageTemperature})
}

// GetGroupSpecies 	godoc
// @Summary      Get the full list of species (with counts) currently detected inside the group
// @Description  Responds with the full list of species (with counts) currently detected inside the group
// @Tags         statistics
// @Produce      json
// Success      200 {object} SpeciesCount
// @Router	    /statistics/group/{groupName}/species [get]
// @Param groupName path string true "Name of the group" example(alpha) default(alpha)
func (sc *StatisticsController) GetGroupSpecies(ctx *gin.Context) {
	groupName := ctx.Param("groupName")

	type SpeciesCount struct {
		Name  string `gorm:"column:name"`
		Count int    `gorm:"column:species_count"`
	}

	var speciesCounts []SpeciesCount
	var sensorObservation models.SensorObservation

	err := sc.DB.Model(&sensorObservation).
		Select("fish_species.name, COUNT(sensor_observations.fish_species_id) AS species_count").
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Joins("JOIN sensor_groups ON sensors.sensor_group_id = sensor_groups.id").
		Joins("JOIN fish_species ON sensor_observations.fish_species_id = fish_species.id").
		Where("sensor_groups.name = ?", groupName).
		Group("fish_species.name").
		Find(&speciesCounts).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve group species"})
		return
	}

	ctx.JSON(http.StatusOK, speciesCounts)
}

// GetRegionMinTemperature godoc
// @Summary      Get the current minimum temperature inside the region
// @Description  Responds with the current minimum temperature inside the region
// @Description  The region is the area represented by the range of coordinates
// @Tags         statistics
// @Produce      json
// @Router /statistics/region/temperature/max [get]
// @Param xMin query int true "Minimum x coordinate of the region" example(1000) default(100)
// @Param xMax query int true "Maximum x coordinate of the region" example(500) default(500)
// @Param yMin query int true "Minimum y coordinate of the region" example(0) default(0)
// @Param yMax query int true "Maximum y coordinate of the region" example(1000) default(1000)
// @Param zMin query int true "Minimum z coordinate of the region" example(0) default(0)
// @Param zMax query int true "Maximum z coordinate of the region" example(5000) default(5000)
func (sc *StatisticsController) GetRegionMinTemperature(xMin, xMax, yMin, yMax, zMin, zMax string) float64 {
	var minTemperature float64
	var sensorObservation models.SensorObservation
	err := sc.DB.Model(&sensorObservation).
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Where("sensors.coordinate_x >= ? AND sensors.coordinate_x <= ?", xMin, xMax).
		Where("sensors.coordinate_y >= ? AND sensors.coordinate_y <= ?", yMin, yMax).
		Where("sensors.coordinate_z >= ? AND sensors.coordinate_z <= ?", zMin, zMax).
		Select("MIN(sensor_observations.temperature) AS minimum_temperature").
		Scan(&minTemperature).Error

	if err != nil {

	}

	return minTemperature
}

// GetRegionMaxTemperature godoc
// @Summary      Get the current maximum temperature inside the region
// @Description  Responds with the current minimum temperature inside the region
// @Description  The region is the area represented by the range of coordinates
// @Tags         statistics
// @Produce      json
// @Router /statistics/region/temperature/min [get]
// @Param xMin query int true "Minimum x coordinate of the region" example(1000) default(100)
// @Param xMax query int true "Maximum x coordinate of the region" example(500) default(500)
// @Param yMin query int true "Minimum y coordinate of the region" example(0) default(0)
// @Param yMax query int true "Maximum y coordinate of the region" example(1000) default(1000)
// @Param zMin query int true "Minimum z coordinate of the region" example(0) default(0)
// @Param zMax query int true "Maximum z coordinate of the region" example(5000) default(5000)
func (sc *StatisticsController) GetRegionMaxTemperature(xMin, xMax, yMin, yMax, zMin, zMax string) float64 {
	var maxTemperature float64
	var sensorObservation models.SensorObservation

	err := sc.DB.Model(&sensorObservation).
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Where("sensors.coordinate_x >= ? AND sensors.coordinate_x <= ?", xMin, xMax).
		Where("sensors.coordinate_y >= ? AND sensors.coordinate_y <= ?", yMin, yMax).
		Where("sensors.coordinate_z >= ? AND sensors.coordinate_z <= ?", zMin, zMax).
		Select("MAX(sensor_observations.temperature) AS maximum_temperature").
		Scan(&maxTemperature).Error

	if err != nil {
	}

	return maxTemperature
}

// GetSensorTemperatureAverageBetweenDates  		   godoc
// @Summary      Get the average temperature detected by a particular sensor between the specified date/time pairs (UNIX timestamps)
// @Description  Responds with the average temperature detected by a particular sensor between the specified date/time pairs (UNIX timestamps)
// @Tags         statistics
// @Produce      json
// @Router /statistics/sensor/{codename}/temperature/average [get]
// @Param codename path string true "Codename of the sensor" example(alpha 5) default(alpha 5)
// @Param from query string true "Start date of the interval" example(1689231951) default(1689231951)
// @Param till query string true "End date of the interval" example(1689262458) default(1689262458)
func (sc *StatisticsController) GetSensorTemperatureAverageBetweenDates(codename string, from, till string) float64 {
	fromInt64, _ := strconv.ParseInt(from, 10, 64)
	tillInt64, _ := strconv.ParseInt(till, 10, 64)
	fromTime := time.Unix(fromInt64, 0)
	tillTime := time.Unix(tillInt64, 0)
	fromTimeString := fromTime.Format("2006-01-02 15:04:05.000000 -07:00")
	tillTimeString := tillTime.Format("2006-01-02 15:04:05.000000 -07:00")
	var averageTemperature sql.NullFloat64

	var sensorObservation models.SensorObservation
	err := sc.DB.Model(&sensorObservation).
		Select("AVG(sensor_observations.temperature) AS average_temperature").
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Where("sensors.codename = ?", codename).
		Where("sensor_observations.detection_time_from >= ?", fromTimeString).
		Where("sensor_observations.detection_time_until <= ?", tillTimeString).
		Scan(&averageTemperature).Error

	if err != nil {
		// Handle the error, e.g., return a default value or error response
	}
	if averageTemperature.Valid == false {
		return -274.0
	} else {
		return averageTemperature.Float64
	}
}

func (sc *StatisticsController) GetTopGroupSpecies(groupName string, N string) []SpeciesCount {
	nLimit, err := strconv.Atoi(N)
	if err != nil {
		return nil
	}

	var speciesCounts []SpeciesCount

	var sensorObservation models.SensorObservation

	err = sc.DB.Model(&sensorObservation).
		Select("fish_species.name, COUNT(sensor_observations.fish_species_id) AS species_count").
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Joins("JOIN sensor_groups ON sensors.sensor_group_id = sensor_groups.id").
		Joins("JOIN fish_species ON sensor_observations.fish_species_id = fish_species.id").
		Where("sensor_groups.name = ?", groupName).
		Group("fish_species.name").
		Order("species_count DESC").
		Limit(nLimit).
		Find(&speciesCounts).Error

	if err != nil {
		return nil
	}

	return speciesCounts
}

// GetTopGroupSpeciesBetweenDates  		   godoc
// @Summary      Get the top N species of a group between two dates
// @Description  Responds with the top N species of a group between two dates
// @Description  Dates are not required
// @Tags         statistics
// @Produce      json
// @Success      200 {object} []SpeciesCount
// @Router /statistics/group/{groupName}/species/top/{N} [get]
// @Param groupName path string true "Name of the group" example(alpha) default(alpha)
// @Param N path int true "Number of species to return" example(5) default(5)
// @Param from query int false "Start date of the interval" example(1689231951) default(1689231951)
// @Param till query int false "End date of the interval" example(1689262458) default(1689262458)
func (sc *StatisticsController) GetTopGroupSpeciesBetweenDates(groupName string, N string, from, till string) []SpeciesCount {

	nLimit, err := strconv.Atoi(N)
	if err != nil {
		return nil
	}
	fromInt64, _ := strconv.ParseInt(from, 10, 64)
	tillInt64, _ := strconv.ParseInt(till, 10, 64)
	fromTime := time.Unix(fromInt64, 0)
	tillTime := time.Unix(tillInt64, 0)
	fromTimeString := fromTime.Format("2006-01-02 15:04:05.000000 -07:00")
	tillTimeString := tillTime.Format("2006-01-02 15:04:05.000000 -07:00")

	var speciesCounts []SpeciesCount

	var sensorObservation models.SensorObservation

	err = sc.DB.Model(&sensorObservation).
		Select("fish_species.name, COUNT(sensor_observations.fish_species_id) AS species_count").
		Joins("JOIN sensors ON sensor_observations.sensor_id = sensors.id").
		Joins("JOIN sensor_groups ON sensors.sensor_group_id = sensor_groups.id").
		Joins("JOIN fish_species ON sensor_observations.fish_species_id = fish_species.id").
		Where("sensor_groups.name = ?", groupName).
		Where("sensor_observations.detection_time_from >= ?", fromTimeString).
		Where("sensor_observations.detection_time_until <= ?", tillTimeString).
		Group("fish_species.name").
		Order("species_count DESC").
		Limit(nLimit).
		Find(&speciesCounts).Error

	if err != nil {
		return nil
	}

	return speciesCounts
}
