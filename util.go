package gatsman

import (
	"errors"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

//RandMap creates random map with specificed size and city count
func RandMap(size int, citycount int) map[int]*City {
	cities := make([]City, size*size)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			cities[(x*size)+y] = City{float64(x), float64(y)}
		}
	}

	cityMap := make(map[int]*City)
	randIndexes := rand.Perm(len(cities) - 1)
	for i := 1; i <= citycount; i++ {
		cityMap[i] = &cities[randIndexes[i]]
	}
	return cityMap
}

//LoadMap function loads the city grid text file and returns a map of city objects
func LoadMap(fileName string) map[int]*City {
	cityMap := make(map[int]*City)

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	s := strings.Split(string(b), "\n")
	info := strings.Split(s[0], " ")
	size, err := strconv.Atoi(strings.Split(info[0], "=")[1])
	if err != nil {
		log.Fatal(err)
	}
	cities, err := strconv.Atoi(strings.Split(info[1], "=")[1])
	if err != nil {
		log.Fatal(err)
	}

	//y should be size - y when creating cities because we start at the top
	for y, line := range s[1:] {
		splitLine := strings.Split(line, " ")
		for x, c := range splitLine {
			if c != "" {
				key, err := strconv.Atoi(c)
				if err != nil {
					log.Print(err)
				} else {
					if key != 0 {
						cityMap[key] = &City{float64(x), float64(size - y - 1)}
					}
				}
			}
		}
	}
	if cities != len(cityMap) {
		log.Fatal(errors.New("error parsing city file"))
	}
	return cityMap
}

//Distance Method to calculate the distance between two cities
func Distance(c1 *City, c2 *City) float64 {
	a := c1.X - c2.X
	b := c1.Y - c2.Y
	return math.Sqrt(a*a + b*b)
}

//IntArrIdxOf finds the index of an int value in an int slice
func IntArrIdxOf(arr []int, value int) int {
	for p, v := range arr {
		if v == value {
			return p
		}
	}
	//log.Print(errors.New("value not in slice"))
	return -1
}

//IntArrCont checks if the given value is in the array or slice
func IntArrCont(arr []int, value int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

//GetNonAdded returns the next int element from one of the two given arrays that is not yet in the third slice
func GetNonAdded(added []int) int {
	shuffled := rand.Perm(CityCount())

	for _, v := range shuffled {
		if !IntArrCont(added, v+1) {
			return v + 1
		}
	}

	return -1
}

//FindString finds the index of a string in an array
func FindString(arr []string, str string) int {
	for p, s := range arr {
		if s == str {
			return p
		}
	}
	//log.Print(errors.New("value not in slice"))
	return -1
}
