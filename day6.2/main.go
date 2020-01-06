package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

func pathToCenterOfMass(start *celestialObject) []*celestialObject {
	path := make([]*celestialObject, 0)

	center := start.center

	for center != nil {
		path = append([]*celestialObject{center}, path...)
		center = center.center
	}

	return path
}

func lowestCommonAncestor(pathA []*celestialObject, pathB []*celestialObject) *celestialObject {
	if len(pathA) == 0 || len(pathB) == 0 || pathA[0] != pathB[0] {
		return nil
	}

	i := 0
	ancestor := pathA[i]

	for i+1 < len(pathA) && i+1 < len(pathB) && pathA[i+1] == pathB[i+1] {
		i++
		ancestor = pathA[i]
	}

	return ancestor
}

func pathLength(from *celestialObject, to *celestialObject) int {
	if from == nil || to == nil {
		return math.MaxInt32
	}

	if from == to {
		return 0
	}

	return 1 + pathLength(from.center, to)
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

	pathFromYou := pathToCenterOfMass(celestialObjects["YOU"])
	pathFromSanta := pathToCenterOfMass(celestialObjects["SAN"])
	ancestor := lowestCommonAncestor(pathFromYou, pathFromSanta)
	lengthPathFromYouToAncestor := pathLength(celestialObjects["YOU"], ancestor)
	lengthPathFromSanToAncestor := pathLength(celestialObjects["SAN"], ancestor)
	fmt.Printf("Ancestor: %s. Orbital Transfers: %d\n", ancestor.name, lengthPathFromYouToAncestor+lengthPathFromSanToAncestor-2)
}
