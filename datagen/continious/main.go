package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/danyaobertan/oceanstats/datagen"
	"github.com/danyaobertan/oceanstats/models"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	epipelagicSpecies    = []string{"Atlantic Bluefin Tuna", "Atlantic Cod", "Atlantic Salmon", "Banded Butterflyfish", "Blue Marlin", "Blue Tang", "Bluebanded Goby", "California Grunion", "Clown Triggerfish", "French Angelfish", "Great Barracuda", "Green Moray Eel", "John Dory", "Mexican Lookdown", "Nassau Grouper", "Northern Red Snapper", "Pacific Halibut", "Pacific Herring", "Pacific Sardine", "Pink Salmon", "Queen Angelfish", "Queen Parrotfish", "Sailfish", "Skipjack Tuna", "Summer Flounder", "Swordfish", "Threespot Damselfish", "Tropical Two-wing Flyingfish", "Wahoo", "Yellowfin Tuna"}
	mesopelagicSpecies   = []string{"Atlantic Wolffish", "Chilean Jack Mackerel", "Chinook Salmon", "Flashlight Fish", "Leafy Seadragon", "Longsnout Seahorse", "Patagonian Toothfish", "Red Lionfish", "Scorpionfish", "Slender Snipe Eel", "Smalltooth Sawfish", "Spotted Moray", "Spotted Porcupinefish", "Spotted Ratfish", "Stoplight Loosejaw", "Tan Bristlemouth", "Whiptail Gulper"}
	bathypelagicSpecies  = []string{"Atlantic Goliath Grouper", "Coelacanth", "Deep Sea Anglerfish", "Oarfish", "Orange Roughy", "Pacific Blackdragon", "Pygmy Seahorse", "Scarlet Frogfish", "Stonefish"}
	abyssopelagicSpecies = []string{"Common Clownfish", "Common Dolphinfish", "Common Fangtooth", "Peruvian Anchoveta"}
	trenchesSpecies      = []string{"Beluga Sturgeon", "Guineafowl Puffer", "Sarcastic Fringehead", "White-ring Garden Eel"}
)

type SensorData struct {
	Data   []models.Sensor `json:"data"`
	Status string          `json:"status"`
}

func main() {
	sensors := GetSensorsFromDB()
	fmt.Println("successfully got sensors from db")
	ContinuousDataGeneration(sensors)
}

func GetSensorsFromDB() []models.Sensor {

	// Send GET request to the API endpoint
	resp, err := http.Get("http://localhost:8000/api/sensor")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// Unmarshal the JSON response into a slice of Sensor structs
	var sensorData SensorData
	err = json.Unmarshal([]byte(body), &sensorData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}
	sensors := sensorData.Data
	return sensors
}

func ContinuousDataGeneration(sensors []models.Sensor) {
	epipelagicSensors := make([]models.Sensor, 0)
	mesopelagicSensors := make([]models.Sensor, 0)
	bathypelagicSensors := make([]models.Sensor, 0)
	abyssopelagicSensors := make([]models.Sensor, 0)
	trenchesSensors := make([]models.Sensor, 0)

	for _, sensor := range sensors {
		_, epipelagicOK := datagen.EpipelagicSensorGroupNames[strings.Split(sensor.Codename, " ")[0]]
		if epipelagicOK {
			epipelagicSensors = append(epipelagicSensors, sensor)
			//fmt.Println("found value:", sensor.Codename)
		}
		_, mesopelagicOK := datagen.MesopelagicSensorGroupNames[strings.Split(sensor.Codename, " ")[0]]
		if mesopelagicOK {
			mesopelagicSensors = append(mesopelagicSensors, sensor)
			//fmt.Println("found value:", sensor.Codename)
		}
		_, bathypelagicOK := datagen.BathypelagicSensorGroupNames[strings.Split(sensor.Codename, " ")[0]]
		if bathypelagicOK {
			bathypelagicSensors = append(bathypelagicSensors, sensor)
			//fmt.Println("found value:", sensor.Codename)
		}
		_, abyssopelagicOK := datagen.AbyssopelagicSensorGroupNames[strings.Split(sensor.Codename, " ")[0]]
		if abyssopelagicOK {
			abyssopelagicSensors = append(abyssopelagicSensors, sensor)
			//fmt.Println("found value:", sensor.Codename)
		}
		_, trenchesOK := datagen.AbyssopelagicSensorGroupNames[strings.Split(sensor.Codename, " ")[0]]
		if trenchesOK {
			trenchesSensors = append(trenchesSensors, sensor)
			//fmt.Println("found value:", sensor.Codename)
		}
	}
	observationsCh := make(chan []models.CreateSensorObservationRequest)

	go GenerateObservations(datagen.EpipelagicDataRate, epipelagicSensors, observationsCh)
	go GenerateObservations(datagen.MesopelagicDataRate, mesopelagicSensors, observationsCh)
	go GenerateObservations(datagen.BathypelagicDataRate, bathypelagicSensors, observationsCh)
	go GenerateObservations(datagen.AbyssopelagicDataRate, abyssopelagicSensors, observationsCh)
	go GenerateObservations(datagen.TrenchesDataRate, trenchesSensors, observationsCh)

	go WriteObservationsToDatabase(observationsCh)

	select {}
}
func WriteObservationsToDatabase(observationsCh <-chan []models.CreateSensorObservationRequest) {
	bufferSize := 1000 // Set the buffer size to control the number of observations to write at once
	observationBuffer := make([]models.CreateSensorObservationRequest, 0, bufferSize)

	// Create a ticker to trigger periodic writes
	ticker := time.NewTicker(10 * time.Minute) // Adjust the duration as needed

	for {
		select {
		case observations := <-observationsCh:
			// Add the received observations to the buffer
			observationBuffer = append(observationBuffer, observations...)

			// If the buffer is full, or if write is triggered by the ticker, write the observations to the database
			if len(observationBuffer) >= bufferSize {
				writeObservationsToDatabase(observationBuffer)
				observationBuffer = observationBuffer[:0] // Clear the buffer
			}

		case <-ticker.C:
			// If a write is triggered by the ticker, write the observations to the database
			writeObservationsToDatabase(observationBuffer)
			observationBuffer = observationBuffer[:0] // Clear the buffer
		}
	}
}

func writeObservationsToDatabase(observations []models.CreateSensorObservationRequest) {
	apiURL := "http://localhost:8000/api/observation/bulk" // API endpoint URL

	// Create the request payload
	payload := struct {
		Observations []models.CreateSensorObservationRequest `json:"observations"`
	}{
		Observations: observations,
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal JSON payload:", err)
		return
	}

	// Send POST request to the API endpoint
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Failed to send HTTP POST request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bulk insertion failed with status:", resp.StatusCode)
		return
	}

	fmt.Println("Bulk insertion completed successfully")
}

func GenerateObservations(datarate int, sensors []models.Sensor, observationsCh chan<- []models.CreateSensorObservationRequest) {
	for {
		genStart := time.Now()

		observations := GenerateObservationsData(datarate, sensors)
		observationsCh <- observations

		genEnd := time.Now()
		genTime := genEnd.Sub(genStart)
		sleepTime := time.Duration(datarate)*time.Second - genTime
		time.Sleep(sleepTime)
	}

}

func GenerateObservationsData(datarate int, sensors []models.Sensor) (observations []models.CreateSensorObservationRequest) {
	observationTime := time.Now()
	observations = make([]models.CreateSensorObservationRequest, 0)
	switch datarate {
	case datagen.EpipelagicDataRate:
		for _, sensor := range sensors {
			temperature := calculateEpipelagicTemperature(sensor.CoordinateZ, observationTime)
			transparency := calculateEpipelagicTransparency(sensor.CoordinateZ)

			fishSpeciesObserved := datagen.GenerateRandomInt(0, len(epipelagicSpecies)/10+2)
			randomIndexes := getRandomUniqueIndexes(len(epipelagicSpecies), fishSpeciesObserved)
			for _, element := range randomIndexes {
				observation := models.CreateSensorObservationRequest{
					SensorID:           sensor.ID,
					Temperature:        temperature,
					Transparency:       int64(transparency),
					FishSpeciesID:      uint(element),
					Count:              datagen.GenerateRandomInt(1, 6),
					DetectionTimeFrom:  observationTime,
					DetectionTimeUntil: observationTime.Add(time.Duration(datagen.EpipelagicDataRate) * time.Second),
				}
				observations = append(observations, observation)
			}
		}
	case datagen.MesopelagicDataRate:
		for _, sensor := range sensors {
			temperature := calculateMesopelagicTemperature(sensor.CoordinateZ, observationTime)
			transparency := calculateMesopelagicZoneTransparency(sensor.CoordinateZ)

			fishSpeciesObserved := datagen.GenerateRandomInt(0, len(mesopelagicSpecies)/10+2)
			randomIndexes := getRandomUniqueIndexes(len(mesopelagicSpecies), fishSpeciesObserved)
			for _, element := range randomIndexes {
				observation := models.CreateSensorObservationRequest{
					SensorID:           sensor.ID,
					Temperature:        temperature,
					Transparency:       int64(transparency),
					FishSpeciesID:      uint(element + len(epipelagicSpecies)),
					Count:              datagen.GenerateRandomInt(1, 6),
					DetectionTimeFrom:  observationTime,
					DetectionTimeUntil: observationTime.Add(time.Duration(datagen.MesopelagicDataRate) * time.Second),
				}
				observations = append(observations, observation)
			}
		}
	case datagen.BathypelagicDataRate:
		for _, sensor := range sensors {
			temperature := calculateBathypelagicTemperature(sensor.CoordinateZ, observationTime)
			transparency := calculateBathypelagicZoneTransparency(sensor.CoordinateZ)

			fishSpeciesObserved := datagen.GenerateRandomInt(0, len(bathypelagicSpecies)/10+2)
			randomIndexes := getRandomUniqueIndexes(len(bathypelagicSpecies), fishSpeciesObserved)
			for _, element := range randomIndexes {
				observation := models.CreateSensorObservationRequest{
					SensorID:           sensor.ID,
					Temperature:        temperature,
					Transparency:       int64(transparency),
					FishSpeciesID:      uint(element + len(epipelagicSpecies) + len(mesopelagicSpecies)),
					Count:              datagen.GenerateRandomInt(1, 6),
					DetectionTimeFrom:  observationTime,
					DetectionTimeUntil: observationTime.Add(time.Duration(datagen.BathypelagicDataRate) * time.Second),
				}
				observations = append(observations, observation)
			}
		}
	case datagen.AbyssopelagicDataRate:
		for _, sensor := range sensors {
			temperature := calculateAbyssopelagicTemperature(sensor.CoordinateZ, observationTime)
			transparency := calculateAbyssopelagicZoneTransparency(sensor.CoordinateZ)

			fishSpeciesObserved := datagen.GenerateRandomInt(0, len(abyssopelagicSpecies)/10+2)
			randomIndexes := getRandomUniqueIndexes(len(abyssopelagicSpecies), fishSpeciesObserved)
			for _, element := range randomIndexes {
				observation := models.CreateSensorObservationRequest{
					SensorID:           sensor.ID,
					Temperature:        temperature,
					Transparency:       int64(transparency),
					FishSpeciesID:      uint(element + len(epipelagicSpecies) + len(mesopelagicSpecies) + len(bathypelagicSpecies)),
					Count:              datagen.GenerateRandomInt(1, 6),
					DetectionTimeFrom:  observationTime,
					DetectionTimeUntil: observationTime.Add(time.Duration(datagen.AbyssopelagicDataRate) * time.Second),
				}
				observations = append(observations, observation)
			}
		}
	case datagen.TrenchesDataRate:
		for _, sensor := range sensors {
			temperature := calculateTrenchesTemperature(sensor.CoordinateZ, observationTime)
			transparency := calculateTrenchesZoneTransparency(sensor.CoordinateZ)

			fishSpeciesObserved := datagen.GenerateRandomInt(0, len(trenchesSpecies)/10+2)
			randomIndexes := getRandomUniqueIndexes(len(trenchesSpecies), fishSpeciesObserved)
			for _, element := range randomIndexes {
				observation := models.CreateSensorObservationRequest{
					SensorID:           sensor.ID,
					Temperature:        temperature,
					Transparency:       int64(transparency),
					FishSpeciesID:      uint(element + len(epipelagicSpecies) + len(mesopelagicSpecies) + len(bathypelagicSpecies) + len(abyssopelagicSpecies)),
					Count:              datagen.GenerateRandomInt(1, 6),
					DetectionTimeFrom:  observationTime,
					DetectionTimeUntil: observationTime.Add(time.Duration(datagen.TrenchesDataRate) * time.Second),
				}
				observations = append(observations, observation)
			}
		}
	}

	return observations
}

// Function to get n random unique indexes from an array
func getRandomUniqueIndexes(length, n int) []int {
	rand.Seed(time.Now().UnixNano())

	// Create a slice to store the selected indexes
	selectedIndexes := make([]int, 0, n)

	if n > length {
		n = length
	}

	// Generate random unique indexes
	for len(selectedIndexes) < n {
		index := rand.Intn(length)
		if !contains(selectedIndexes, index) {
			selectedIndexes = append(selectedIndexes, index)
		}
	}

	return selectedIndexes
}

// Helper function to check if a slice contains a specific element
func contains(slice []int, element int) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

// Helper functions to calculate temperature profiles in different ocean zones
func calculateEpipelagicTemperature(depth float64, timeOfDay time.Time) (temperature float64) {
	// Calculate temperature based on depth and time of the day in the Epipelagic zone
	meanTemperature := 25.0 // Mean temperature in °C
	amplitude := 2.0
	periodInHours := 24.0

	timeInHours := float64(timeOfDay.Hour()) // Convert current time in hours

	temperature = meanTemperature + amplitude*math.Sin(2*math.Pi*timeInHours/periodInHours+math.Pi/2)

	//fmt.Printf("Current temperature: %.2f°C\n", temperature)

	attenuationRate := 0.020
	temperature -= attenuationRate * depth
	//fmt.Printf("After attenuationRate temperature: %.2f°C\n", temperature)

	return math.Round(temperature*100) / 100
}

// Calculate temperature based on depth and time of the day in the Twilight zone
func calculateMesopelagicTemperature(depth float64, timeOfDay time.Time) float64 {
	meanTemperature := 18.0
	amplitude := 1.0
	periodInHours := 24.0

	timeInHours := float64(timeOfDay.Hour())

	temperature := meanTemperature + amplitude*math.Sin(2*math.Pi*timeInHours/periodInHours+math.Pi/2)

	//fmt.Printf("Current temperature in the Twilight zone: %.2f°C\n", temperature)

	attenuationRate := 0.015 / depth
	temperature -= attenuationRate * depth
	//fmt.Printf("After attenuation rate in the Twilight zone: %.2f°C\n", temperature)

	return math.Round(temperature*100) / 100
}

// Calculate temperature based on depth and time of the day in the Midnight zone
func calculateBathypelagicTemperature(depth float64, timeOfDay time.Time) float64 {
	meanTemperature := 15.0
	amplitude := 0.5
	periodInHours := 24.0

	timeInHours := float64(timeOfDay.Hour())

	temperature := meanTemperature + amplitude*math.Sin(2*math.Pi*timeInHours/periodInHours+math.Pi/2)

	//fmt.Printf("Current temperature in the Midnight zone: %.2f°C\n", temperature)

	attenuationRate := 0.010 / depth
	temperature -= attenuationRate * depth
	//fmt.Printf("After attenuation rate in the Midnight zone: %.2f°C\n", temperature)

	return math.Round(temperature*100) / 100
}

// Calculate temperature based on depth and time of the day in the Abyssal zone
func calculateAbyssopelagicTemperature(depth float64, timeOfDay time.Time) float64 {
	meanTemperature := 12.0
	amplitude := 0.2
	periodInHours := 24.0

	timeInHours := float64(timeOfDay.Hour())

	temperature := meanTemperature + amplitude*math.Sin(2*math.Pi*timeInHours/periodInHours+math.Pi/2)

	//fmt.Printf("Current temperature in the Abyssal zone: %.2f°C\n", temperature)

	attenuationRate := 0.005 / depth
	temperature -= attenuationRate * depth
	//fmt.Printf("After attenuation rate in the Abyssal zone: %.2f°C\n", temperature)

	return math.Round(temperature*100) / 100
}

// Calculate temperature based on depth and time of the day in the Hadal zone
func calculateTrenchesTemperature(depth float64, timeOfDay time.Time) float64 {
	meanTemperature := 10.0
	amplitude := 0.1
	periodInHours := 24.0

	timeInHours := float64(timeOfDay.Hour())

	temperature := meanTemperature + amplitude*math.Sin(2*math.Pi*timeInHours/periodInHours+math.Pi/2)

	//fmt.Printf("Current temperature in the Hadal zone: %.2f°C\n", temperature)

	attenuationRate := 0.002 / depth
	temperature -= attenuationRate * depth
	//fmt.Printf("After attenuation rate in the Hadal zone: %.2f°C\n", temperature)

	return math.Round(temperature*100) / 100
}

// Helper functions to calculate water transparency percentages in different depth zones
func calculateEpipelagicTransparency(depth float64) int {
	// Calculate water transparency as a percentage for the Epipelagic zone using a formula
	transparency := 100.0 - (depth/10.0)*2.7
	if transparency < 0.0 {
		transparency = 0.0
	}
	return int(int64(transparency + datagen.GenerateRandomFloat(-3.0, 3.0)))
}

func calculateMesopelagicZoneTransparency(depth float64) int {
	// Calculate water transparency as a percentage for the Twilight Zone using a formula
	transparency := 50.0 - (depth-10.0)/190.0*4.0
	if transparency < 0.0 {
		transparency = 0.0
	}
	return int(transparency + datagen.GenerateRandomFloat(-2.0, 2.0))
}

func calculateBathypelagicZoneTransparency(depth float64) int {
	// Calculate water transparency as a percentage for the Midnight Zone using a formula
	transparency := 30.0 - (depth-200.0)/1800.0*11.0
	if transparency < 0.0 {
		transparency = 0.0
	}
	return int(transparency + datagen.GenerateRandomFloat(-1.0, 1.0))
}

func calculateAbyssopelagicZoneTransparency(depth float64) int {
	// Calculate water transparency as a percentage for the Abyssal Zone using a formula
	transparency := 10.0 - (depth-1000.0)/3000.0*4.0
	if transparency < 0.0 {
		transparency = 0.0
	}
	return int(transparency + datagen.GenerateRandomFloat(-1.0, 1.0))
}

func calculateTrenchesZoneTransparency(depth float64) int {
	// Calculate water transparency as a percentage for the Hadal Zone using a formula
	transparency := 5.0 - (depth-4000.0)/6000.0*6.0
	if transparency < 0.0 {
		transparency = 0.0
	}
	return int(transparency + datagen.GenerateRandomFloat(-0.5, 0.5))
}
