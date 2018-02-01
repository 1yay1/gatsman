package gatsman

import (
	"errors"
	"log"
)

//CityPair is a pair of two city indexes in the citymap.
type CityPair struct {
	C1 int
	C2 int
}

//NewCityPair creates a new citypair object
func NewCityPair(c1 int, c2 int) CityPair {
	if c1 == c2 {
		log.Fatal(errors.New("can't create CityPair from equal cities"))
	}
	if c1 > c2 {
		return CityPair{c1, c2}
	}
	return CityPair{c2, c1}

}
