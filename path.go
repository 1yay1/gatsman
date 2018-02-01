package gatsman

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//NewRandomPath creates a new slice of unique ints, which is just a permutation of all keys in the cityMap
func NewRandomPath() []int {

	p := rand.Perm(CityCount())
	for i := 0; i < len(p); i++ {
		p[i]++
	}

	return p
}
