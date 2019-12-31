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

		center := relationshipCelestialObjects[0]
		satellite := relationshipCelestialObjects[1]

		var centerObject, satelliteObject celestialObject
		if centerObject, ok := celestialObjects[center]; ok {
			if satelliteObject, ok := celestialObjects[satellite]; ok {
				satelliteObject.center = centerObject
			} else {
				satelliteObject = &celestialObject{name: satellite, center: centerObject, satellites: make([]*celestialObject, 0)}
				celestialObjects[satellite] = satelliteObject
			}
		} else {
			centerObject := &celestialObject{name: center, satellites: make([]*celestialObject, 0)}

			if satelliteObject, ok := celestialObjects[satellite]; ok {
				centerObject.satellites = append(centerObject.satellites, satelliteObject)
				satelliteObject.center = centerObject
			} else {
				satelliteObject = &celestialObject{name: satellite, center: centerObject, satellites: make([]*celestialObject, 0)}
				celestialObjects[satellite] = satelliteObject
			}

			celestialObjects[center] = centerObject
		}

		centerObject.satellites = append(centerObject.satellites, &satelliteObject)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	celestialObjects := make(map[string]*celestialObject)

	err := parseOrbitalRelationships("test.txt", celestialObjects)

	if err != nil {
		panic(err)
	}

	for name, celestialObject := range celestialObjects {
		fmt.Printf("Celestial Object %s with %d satellites\n", name, len(celestialObject.satellites))
	}
}
