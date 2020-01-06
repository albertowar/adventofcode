package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type celestialObject struct {
	name       string
	center     *celestialObject
	satellites []*celestialObject
}

func parseOrbitalRelationships(filename string, celestialObjects map[string]*celestialObject) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		orbitalRelationship := scanner.Text()

		relationshipCelestialObjects := strings.Split(orbitalRelationship, ")")

		if len(relationshipCelestialObjects) != 2 {
			return fmt.Errorf("Too many objects in the orbit. %s with length %d. Original %s", relationshipCelestialObjects, len(relationshipCelestialObjects), orbitalRelationship)
		}

		centerName := relationshipCelestialObjects[0]
		satelliteName := relationshipCelestialObjects[1]

		centerObject, found := celestialObjects[centerName]
		if !found {
			centerObject = &celestialObject{name: centerName, satellites: make([]*celestialObject, 0)}
			celestialObjects[centerName] = centerObject
		}

		satelliteObject, found := celestialObjects[satelliteName]
		if found {
			centerObject.satellites = append(centerObject.satellites, satelliteObject)
			satelliteObject.center = centerObject
		} else {
			satelliteObject = &celestialObject{name: satelliteName, center: centerObject, satellites: make([]*celestialObject, 0)}
			celestialObjects[satelliteName] = satelliteObject
		}

		centerObject.satellites = append(centerObject.satellites, satelliteObject)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func orbits(current *celestialObject, center *celestialObject) int {
	if center == nil {
		return 0
	}

	return 1 + orbits(current, center.center)
}

func main() {
	celestialObjects := make(map[string]*celestialObject)

	err := parseOrbitalRelationships("input.txt", celestialObjects)

	if err != nil {
		panic(err)
	}

	for name, celestialObject := range celestialObjects {
		fmt.Printf("Celestial Object %s with %d satellites\n", name, len(celestialObject.satellites))
	}

	totalOrbits := 0

	for _, celestialObject := range celestialObjects {
		totalOrbits += orbits(celestialObject, celestialObject.center)
	}

	fmt.Printf("Total orbits: %d\n", totalOrbits)
}
