package gatsman

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//Gatsworker that does all the
type Gatsworker struct {
	wg     *sync.WaitGroup
	popc   chan Popu
	outpc  chan OutpData
	mgen   int
	runs   int
	maxFit float64
}

//OutpData out data of the worker
type OutpData struct {
	Pop      Popu
	Gens     float64
	Fit      float64
	Duration int64
}

//ToStringArr returns a string array to be printed as csv
func (od OutpData) ToStringArr() []string {
	return []string{fmt.Sprintf("%v", od.Pop.Pm), fmt.Sprintf("%v", od.Pop.Pc), fmt.Sprintf("%v", od.Gens)}
}

//ToStringArrRand return rand
func (od OutpData) ToStringArrRand() []string {
	return []string{fmt.Sprintf("%v", od.Pop.Pm), fmt.Sprintf("%v", od.Pop.Pc), fmt.Sprintf("%v", od.Fit), fmt.Sprintf("%v", od.Duration)}
}

//NewGatsworker constrcutor for gatsworker
func NewGatsworker(wg *sync.WaitGroup, popc chan Popu, outpc chan OutpData, mgen, runs int, maxFit float64) Gatsworker {
	wg.Add(1)
	return Gatsworker{wg, popc, outpc, mgen, runs, maxFit}
}

//WorkRandom works on map with unknown maxfit
func (gw Gatsworker) WorkRandom(debug bool) {
	rand.Seed(time.Now().UnixNano())
	if debug {
		fmt.Printf("Starting work!\n")
	}
	defer gw.wg.Done()
	for {

		pop, ok := <-gw.popc
		if !ok {
			//nothing left to grab
			return
		}

		for r := 0; r < gw.runs; r++ {
			g := 0
			start := time.Now().UnixNano()
			for ; g < gw.mgen; g++ {
				pop.Evolve()
			}
			end := time.Now().UnixNano()
			if debug {
				fmt.Printf("pm: %v, pc: %v, gens: %v, fit: %v time: %v \n", pop.Pm, pop.Pc, g, pop.Pop[0].Fitness, end-start)
			}
			gw.outpc <- OutpData{pop, float64(g), pop.Pop[0].Fitness, end - start}
			//fmt.Printf("gen: %v gensum: %v fit: %v\n ", float64(g), gens/float64(r+1), pop.Pop[0].Fitness)
			pop = pop.NewCopy()
		}
	}
}

//Work loop of a worker. Grabs inp data from channel aslong as it is filled, works and adds the result to outpc then closes
func (gw Gatsworker) Work(debug bool) {
	rand.Seed(time.Now().UnixNano())
	if debug {
		fmt.Printf("Starting work!\n")
	}
	defer gw.wg.Done()
	for {

		start := time.Now().UnixNano()
		pop, ok := <-gw.popc
		if !ok {
			//nothing left to grab
			return
		}
		fit := 0.0
		gens := 0.0

		for r := 0; r < gw.runs; r++ {
			g := 0
			for ; g < gw.mgen; g++ {
				pop.Evolve()
				if pop.Pop[0].Fitness <= gw.maxFit {
					break
				}
			}
			fit += pop.Pop[0].Fitness
			gens += float64(g)
			//fmt.Printf("gen: %v gensum: %v fit: %v\n ", float64(g), gens/float64(r+1), pop.Pop[0].Fitness)
			pop = pop.NewCopy()
		}

		end := time.Now().UnixNano()
		if debug {
			fmt.Printf("pm: %v, pc: %v, gens: %v, fit: %v time: %v \n", pop.Pm, pop.Pc, gens/float64(gw.runs), fit/float64(gw.runs), end-start)
		}
		gw.outpc <- OutpData{pop, gens / float64(gw.runs), fit / float64(gw.runs), end - start}
	}
}
