package datagen

import (
	"fmt"
	"math/rand"
	"strings"
)

const precision = 10
const epipelagicBorder float64 = 200
const mesopelagicBorder float64 = 1000
const bathypelagicBorder float64 = 4000
const abyssopelagicBorder float64 = 6000

const EpipelagicDataRate int = 15
const MesopelagicDataRate int = 30
const BathypelagicDataRate int = 60
const AbyssopelagicDataRate int = 90
const TrenchesDataRate int = 120

var (
	EpipelagicSensorGroupNames = map[string]bool{
		"alpha":   true,
		"beta":    true,
		"gamma":   true,
		"delta":   true,
		"epsilon": true,
	}

	MesopelagicSensorGroupNames = map[string]bool{
		"zeta":  true,
		"eta":   true,
		"theta": true,
		"iota":  true,
		"kappa": true,
	}

	BathypelagicSensorGroupNames = map[string]bool{
		"lambda":  true,
		"mu":      true,
		"nu":      true,
		"xi":      true,
		"omicron": true,
	}

	AbyssopelagicSensorGroupNames = map[string]bool{
		"pi":      true,
		"rho":     true,
		"sigma":   true,
		"tau":     true,
		"upsilon": true,
	}

	TrenchesSensorGroupNames = map[string]bool{
		"phi":   true,
		"chi":   true,
		"psi":   true,
		"omega": true,
	}
)

type Sensor struct {
	Codename    string
	Coordinates Coordinates
	DataRate    int
}

type Coordinates struct {
	X, Y, Z float64
}

type SensorGroup struct {
	Name    string
	Sensors []Sensor
}

func GenerateSensors() []SensorGroup {
	sensorGroups := []SensorGroup{
		{Name: "alpha"},
		{Name: "beta"},
		{Name: "gamma"},
		{Name: "delta"},
		{Name: "epsilon"},
		{Name: "zeta"},
		{Name: "eta"},
		{Name: "theta"},
		{Name: "iota"},
		{Name: "kappa"},
		{Name: "lambda"},
		{Name: "mu"},
		{Name: "nu"},
		{Name: "xi"},
		{Name: "omicron"},
		{Name: "pi"},
		{Name: "rho"},
		{Name: "sigma"},
		{Name: "tau"},
		{Name: "upsilon"},
		{Name: "phi"},
		{Name: "chi"},
		{Name: "psi"},
		{Name: "omega"},
	}

	// Generating sensors within each group
	for i, group := range sensorGroups {
		//s represents the sensor group index
		s := i % 5
		fmt.Println(i, s)
		//var epipelagicBorder float64 = 200
		//var mesopelagicBorder float64 = 1000
		//var bathypelagicBorder float64 = 4000
		//var abyssopelagicBorder float64 = 6000
		// Generating sensors for Sunlight Zone (Epipelagic)
		// Depth Range: 0-200 meters
		// Length from Shore: Extends from the coast to about 200 meters
		if EpipelagicSensorGroupNames[group.Name] {
			//fmt.Println(i, group.Name)
			for j := 0; j <= 10; j++ {
				sensor := Sensor{
					Codename:    fmt.Sprintf("%s %d", group.Name, j),
					Coordinates: generateCoordinates(float64(j*100), float64(s*40), float64((s+1)*40), float64(s*40), float64((s+1)*40)),
					DataRate:    EpipelagicDataRate,
				}
				sensorGroups[i].Sensors = append(sensorGroups[i].Sensors, sensor)
			}
		}

		// Generating sensors for Twilight Zone (Mesopelagic)
		// Depth Range: 200-1000 meters
		// Length from Shore: Extends from 200 meters to about 1,000 meters
		if MesopelagicSensorGroupNames[group.Name] {
			//fmt.Println(i, group.Name)
			for j := 0; j <= 10; j++ {
				sensor := Sensor{
					Codename:    fmt.Sprintf("%s %d", group.Name, j),
					Coordinates: generateCoordinates(float64(j*200), float64(s*160)+epipelagicBorder, float64((s+1)*160), float64(s*160)+epipelagicBorder, float64((s+1)*160)),
					DataRate:    MesopelagicDataRate,
				}
				sensorGroups[i].Sensors = append(sensorGroups[i].Sensors, sensor)
			}
		}

		// Generating sensors for Midnight Zone (Bathypelagic)
		// Depth Range: 1000-4000 meters
		// Length from Shore: Extends from 1,000 meters to about 4,000 meters
		if BathypelagicSensorGroupNames[group.Name] {
			//fmt.Println(i, group.Name)
			for j := 0; j <= 10; j++ {
				sensor := Sensor{
					Codename:    fmt.Sprintf("%s %d", group.Name, j),
					Coordinates: generateCoordinates(float64(j*400), float64(s*600)+mesopelagicBorder, float64((s+1)*600), float64(s*600)+mesopelagicBorder, float64((s+1)*600)),
					DataRate:    BathypelagicDataRate,
				}
				sensorGroups[i].Sensors = append(sensorGroups[i].Sensors, sensor)
			}
		}

		// Generating sensors for Abyss (Abyssopelagic)
		// Depth Range: 4000-6000 meters
		// Length from Shore: Extends from 4,000 meters to about 6,000 meters
		if AbyssopelagicSensorGroupNames[group.Name] {
			//fmt.Println(i, group.Name)
			for j := 0; j <= 10; j++ {
				sensor := Sensor{
					Codename:    fmt.Sprintf("%s %d", group.Name, j),
					Coordinates: generateCoordinates(float64(j*800), float64(s*400)+bathypelagicBorder, float64((s+1)*400), float64(s*400)+bathypelagicBorder, float64((s+1)*400)),
					DataRate:    AbyssopelagicDataRate,
				}
				sensorGroups[i].Sensors = append(sensorGroups[i].Sensors, sensor)
			}
		}

		// Generating sensors for Trenches
		// Depth Range: 6000-11000 meters
		// Length from Shore: Extends from 6,000 meters to the bottom of the ocean
		if TrenchesSensorGroupNames[group.Name] {
			//fmt.Println(i, group.Name)
			for j := 0; j <= 10; j++ {
				sensor := Sensor{
					Codename:    fmt.Sprintf("%s %d", group.Name, j),
					Coordinates: generateCoordinates(float64(j*1600), float64(s*1250)+abyssopelagicBorder, float64((s+1)*1250), float64(s*1250)+abyssopelagicBorder, float64((s+1)*1250)),
					DataRate:    TrenchesDataRate,
				}
				sensorGroups[i].Sensors = append(sensorGroups[i].Sensors, sensor)
			}
		}

	}

	return sensorGroups
}

func generateCoordinates(X, minY, maxY, minZ, maxZ float64) Coordinates {
	// Generate random coordinates within the desired range
	return Coordinates{
		X: X + generateDisplacement(),
		Y: GenerateRandomFloat(minY, maxY),
		Z: GenerateRandomFloat(minZ, maxZ),
	}
}

//func generateDataRate(min, max int) int {
//	// Generating a random data rate in seconds
//	return generateRandomInt(min, max)
//}

func GenerateRandomFloat(min, max float64) float64 {
	// Generating a random floating-point number within the desired range
	return min + rand.Float64()*(max-min)
}

func GenerateRandomInt(min, max int) int {
	// Generating a random integer within the desired range
	return min + rand.Intn(max-min+1)
}

func generateDisplacement() float64 {
	return GenerateRandomFloat(-precision, precision)
}

func GenerateSQLStatements(sensorGroups []SensorGroup) string {

	var sqlDataGenerationStatements strings.Builder
	sqlDataGenerationStatements.WriteString("INSERT INTO sensors (codename, coordinate_x, coordinate_y, coordinate_z, data_rate, sensor_group_id) \n VALUES")

	for i := range sensorGroups {
		sensors := sensorGroups[i].Sensors
		for j, sensor := range sensors {
			sqlDataGenerationStatements.WriteString(fmt.Sprintf("('%s', %.2f, %.2f, %.2f, %d, %d)", sensor.Codename, sensor.Coordinates.X, sensor.Coordinates.Y, sensor.Coordinates.Z, sensor.DataRate, i+1))

			// Append comma if it is not the last sensor
			if j != len(sensors)-1 || i != len(sensorGroups)-1 {
				sqlDataGenerationStatements.WriteString(",\n")
			}
		}
	}

	// Replace the last comma with a semicolon
	sqlDataGenerationStatements.WriteString(";")

	return sqlDataGenerationStatements.String()
}
