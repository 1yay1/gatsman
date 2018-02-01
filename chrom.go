package gatsman

import "math/rand"

//Tsman stores a path and it's fitness. There are functions to calculate them
type Tsman struct {
	Path    []int
	Fitness float64
}

//CalcFit calculates the complete distance of the path, including a trip from the last city to the first, making it a round trip
func CalcFit(path []int) float64 {
	var fit float64
	for i := 0; i < len(path)-1; i++ {
		cp := NewCityPair(path[i], path[i+1])
		fit += GetDistance(cp)
	}
	cp := NewCityPair(path[0], path[len(path)-1])
	fit += GetDistance(cp)
	return fit
}

//NewTsman returns a new Tsman struct
func NewTsman() *Tsman {
	return NewTsmanFromPath(NewRandomPath())
}

//NewTsmanFromPath returns a new Tsman struct
func NewTsmanFromPath(path []int) *Tsman {
	ts := Tsman{path, CalcFit(path)}
	return &ts
}

//Mutate swaps two random cities on average n(n=len(path))*pm times, except the first one
func (ts *Tsman) Mutate() {
	idx1 := rand.Intn(len(ts.Path) - 1)
	idx2 := rand.Intn(len(ts.Path) - 1)
	ts.Path[idx1], ts.Path[idx2] = ts.Path[idx2], ts.Path[idx1]
	/*if ts.Fitness < CalcFit(ts.Path) {
		ts.Path[idx1], ts.Path[idx2] = ts.Path[idx2], ts.Path[idx1]
	}*/

	//ts.Fitness = CalcFit(ts.Path)
}

//GreedyMutate mutates only if the distances are shorter
func (ts *Tsman) GreedyMutate() {
	idx1 := rand.Intn(len(ts.Path) - 1)
	idx2 := rand.Intn(len(ts.Path) - 1)
	ts.Path[idx1], ts.Path[idx2] = ts.Path[idx2], ts.Path[idx1]
	newFit := CalcFit(ts.Path)
	if ts.Fitness < newFit {
		ts.Path[idx1], ts.Path[idx2] = ts.Path[idx2], ts.Path[idx1]
	} else {
		ts.Fitness = newFit
	}
	//ts.Fitness = CalcFit(ts.Path)
}

//GreddyQuadCross works like greedycross, but checks previous and next city, and random city.CityPair
func GreddyQuadCross(ts1 *Tsman, ts2 *Tsman) *Tsman {
	var newPath = make([]int, CityCount())
	var newFitness float64
	newPath[0] = ts1.Path[0]

	length := len(newPath)
	nextFour := make([]int, 4)
	//for i := 1; i < len(newPath)-1; i++ {
	//get index of previously added to new path in both old paths
	i := 1
	for j := 0; i < length; j++ {
		var added = newPath[0:i]

		idx1 := IntArrIdxOf(ts1.Path, newPath[i-1])
		if idx1 < length-1 {
			nextFour[0] = ts1.Path[idx1+1]
			if idx1 == 0 {
				nextFour[1] = ts1.Path[length-1]
			} else {
				nextFour[1] = ts1.Path[idx1-1]
			}
		} else {
			nextFour[0] = ts1.Path[0]
			nextFour[1] = ts1.Path[length-2]
		}

		idx2 := IntArrIdxOf(ts2.Path, newPath[i-1])
		if idx2 < length-1 {
			nextFour[2] = ts2.Path[idx2+1]
			if idx2 == 0 {
				nextFour[3] = ts2.Path[length-1]
			} else {
				nextFour[3] = ts2.Path[idx2-1]
			}
		} else {
			nextFour[2] = ts2.Path[0]
			nextFour[3] = ts2.Path[length-2]
		}

		for k, v := range nextFour {
			if IntArrCont(added, v) {
				nextFour[k] = GetNonAdded(added)
			}
		}

		next := nextFour[0]
		d := GetDistance(NewCityPair(newPath[i-1], nextFour[0]))
		for k := 1; k < 4; k++ {
			temp := GetDistance(NewCityPair(newPath[i-1], nextFour[k]))
			if temp < d {
				//fmt.Printf("nextFour[%v} = %v, d = %v, temp = %v\n", k, nextFour[k], d, temp)
				next = nextFour[k]
			}
		}
		newPath[i] = next
		newFitness += d
		i++
	}
	newFitness += GetDistance(NewCityPair(newPath[0], newPath[len(newPath)-1]))
	return &Tsman{newPath, newFitness}
}

//GreedyCross returns new Tsman pointer with the newly crossed path, greedy crossover
func GreedyCross(ts1 *Tsman, ts2 *Tsman) *Tsman {
	var newPath = make([]int, CityCount())
	var newFitness float64
	newPath[0] = ts1.Path[0]

	//for i := 1; i < len(newPath)-1; i++ {
	//get index of previously added to new path in both old paths
	i := 1
	for j := 0; j < len(newPath)-1; j++ {
		var added = newPath[0:i]

		var idx1, idx2, from1, to1, from2, to2 int
		var contains1, contains2 bool
		var cp1, cp2 CityPair
		var d1, d2 float64

		idx1 = IntArrIdxOf(ts1.Path, newPath[i-1])
		from1 = ts1.Path[idx1]
		if idx1 < len(newPath)-1 {
			to1 = ts1.Path[idx1+1]
		} else {
			to1 = ts1.Path[0]
		}
		contains1 = IntArrCont(added, to1)

		idx2 = IntArrIdxOf(ts2.Path, newPath[i-1])
		from2 = ts2.Path[idx2]
		if idx2 < len(newPath)-1 {
			to2 = ts2.Path[idx2+1]
		} else {
			to2 = ts2.Path[0]
		}
		contains2 = IntArrCont(added, to2)

		//log.Printf("\ni: %v added: %v\n", i, added)
		//log.Printf("idx1: %v from1: %v to1: %v contains1: %v\n", idx1, from1, to1, contains1)
		//log.Printf("idx2: %v from2: %v to2: %v contains2: %v\n", idx2, from2, to2, contains2)
		if !contains1 && !contains2 {
			//log.Print("Both not added\n")
			cp1 = NewCityPair(from1, to1)
			cp2 = NewCityPair(from2, to2)

			d1 = GetDistance(cp1)
			d2 = GetDistance(cp2)
			//both not added yet
			//minimize distance
			if d1 < d2 {
				newPath[i] = to1
				newFitness += d1
			} else {
				newPath[i] = to2
				newFitness += d2
			}
		} else if contains1 && !contains2 {
			//log.Print("Added first\n")
			//added only first
			newPath[i] = to2
			newFitness += GetDistance(NewCityPair(from2, to2))

		} else if !contains1 && contains2 {
			//log.Print("Added second\n")
			//added only first
			newPath[i] = to1
			newFitness += GetDistance(NewCityPair(from1, to1))
		} else if contains1 && contains2 {
			//log.Print("Added both!")
			//both already added, choose a non-added
			v := GetNonAdded(added)
			if v > 0 {
				newPath[i] = v
				newFitness += GetDistance(NewCityPair(newPath[i-1], newPath[i]))
			} else {
				break
			}
		}

		i++
	}
	newFitness += GetDistance(NewCityPair(newPath[0], newPath[len(newPath)-1]))
	return &Tsman{newPath, newFitness}
}
