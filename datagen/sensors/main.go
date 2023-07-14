package main

import (
	"fmt"
	"github.com/danyaobertan/oceanstats/datagen"
	"io/ioutil"
)

func main() {
	// Generating the sensor groups and sensors
	sensorGroups := datagen.GenerateSensors()
	fmt.Println(sensorGroups)
	// Generating the SQL statements
	sqlDataGenerationStatements := datagen.GenerateSQLStatements(sensorGroups)
	err := ioutil.WriteFile("initialSQL.txt", []byte(sqlDataGenerationStatements), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("String successfully written to file.")
}
