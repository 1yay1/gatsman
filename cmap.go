package gatsman

import "fmt"

//cityMap that stores the loaded cities from file or ranomdly generated
var cityMap map[int]*City

//distanceMap maps a pair of city indexes to the distance between them
var distanceMap map[CityPair]float64

//GetDistance  If the distance has been calculated it is retrieved, else it is added and then returned
func GetDistance(cp CityPair) float64 {
	val, ok := distanceMap[cp]
	if !ok {
		fmt.Printf("error\n")
	}
	return val
}

//LenDistanceMap returns length of distance map
func LenDistanceMap() int {
	return len(distanceMap)
}

func init() {
	//CityMap = make(map[int]*City)

}

//RandInitCMap initializes the CMap with a random map defined by a size and citycount
func RandInitCMap(size int, citycount int) {
	cityMap = RandMap(size, citycount)
}

//InitCMap initializes the cityMap with a map from file
func InitCMap(filename string) {
	if filename == "" {
		RandInitCMap(1000, 50)
	} else {
		cityMap = LoadMap(filename)
	}
}

//InitDMap Initializes the DMap
func InitDMap() {
	distanceMap = make(map[CityPair]float64)
	for i := 1; i <= CityCount(); i++ {
		for j := i + 1; j <= CityCount(); j++ {
			d := Distance(cityMap[i], cityMap[j])
			distanceMap[NewCityPair(i, j)] = d
		}
	}
}

//GetCity function to return the city with the given index
func GetCity(key int) *City {
	return cityMap[key]
}

//CityCount returns the amount of cities
func CityCount() int {
	return len(cityMap)
}
