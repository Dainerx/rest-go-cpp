package vrpreader

import (
	"fmt"
	"log"
	"strings"
)

// Read an instance's input passed as argument.
// Returns lines read and nil if all good, otherwise 1 and an error.
func ReadInstance(input string) (int, error) {
	var trucks, clients, dimension, sum int
	// Split input into different lines.
	lines := strings.Split(input, "\n")
	line := 0
	// Read header.
	_, err := fmt.Sscanf(lines[line], "%d %d %d %d", &clients, &trucks, &dimension, &sum)
	// Checks if the line respects header agreed format, panics if Sscanf fails.
	if err != nil {
		log.Panicf("fmt.Sscanf(input, &clients,&trucks, &dimension, &sum) failed: %v", err)
		return -1, err
	}

	// Read every truck's data.
	for i := 0; i < trucks; i++ {
		var idTruck, capacity, startTime, endTime, idStartPoint, idEndPoint int16
		var latitudeStartPoint, longitudeStartPoint, latitudeEndPoint, longitudeEndPoint float32
		line++
		_, err := fmt.Sscanf(lines[line], "%d %d %d %d %d %f %f %d %f %f", &idTruck, &capacity, &startTime,
			&endTime, &idStartPoint, &latitudeStartPoint, &longitudeStartPoint, &idEndPoint,
			&latitudeEndPoint, &longitudeEndPoint)
		// Checks if the line respects truck agreed format, panics if Sscanf fails.
		if err != nil {
			log.Panicf("Reading truck data with id %d on line %d failed: %v", i, line+1, err) //line + 1 to indicate which line on the file
			return -1, err
		}
	}

	// Read every client's data.
	for i := 0; i < clients; i++ {
		var idClient, demand, serviceTime, startTime, endTime, idPoint int16
		var latitudePoint, longitudePoint float32
		line++
		_, err := fmt.Sscanf(lines[line], "%d %d %d %d %d %d %f %f", &idClient, &demand, &serviceTime,
			&startTime, &endTime, &idPoint, &latitudePoint, &longitudePoint)
		// Checks if the line respects client agreed format, panics if Sscanf fails.
		if err != nil {
			log.Panicf("Reading client data with id %d on line %d failed: %v", i, line+1, err) //line + 1 to indicate which line on the file
			return -1, err
		}
	}

	// All good return lines read and nil as error.
	return line, nil
}
