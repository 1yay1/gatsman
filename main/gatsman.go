package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/1yay1/gatsman"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

func init() {

	//	initDefault()
}

func main() {
	runRandom()
}

func run() {
	psPtr := flag.Int("ps", 100, "population size, integer")
	crossPtr := flag.Int("cross", 0, "crossover method")
	replicPtr := flag.Int("replc", 0, "replication method")
	protectPtr := flag.Bool("protect", true, "protect best bool")
	runPtr := flag.Int("runs", 50, "runs")
	maxFitPtr := flag.Float64("maxfit", 36.0, "max fit breakpoint")
	maxGenPtr := flag.Int("mgen", 1000, "max gens")
	inptFilePtr := flag.String("inpt", "", "inpt file")
	outpFilePtr := flag.String("outpt", fmt.Sprintf("out/outp_%v.csv", time.Now().UnixNano()), "output file")
	debugPtr := flag.Bool("debug", true, "debug flag")
	workerFlg := flag.Int("worker", runtime.NumCPU(), "number of workers to run")
	flag.Parse()

	fmt.Printf(*inptFilePtr)
	inptFile := *inptFilePtr
	gatsman.InitCMap(inptFile)
	gatsman.InitDMap()

	outpFile := *outpFilePtr
	file, err := os.Create(outpFile)
	if err != nil {
		log.Fatal("cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	ps := *psPtr
	cross := *crossPtr
	replic := *replicPtr
	protect := *protectPtr
	mgen := *maxGenPtr
	runs := *runPtr
	maxFit := *maxFitPtr
	debug := *debugPtr
	workers := *workerFlg
	var wg sync.WaitGroup
	popc := make(chan gatsman.Popu, 10000)
	for _, pm := range createSteps(0.0, 0.005, 40) {
		for _, pc := range createSteps(0.0, 0.05, 18) {
			/*
				type Gatsworker struct {
					wg     *sync.WaitGroup
					popc   chan Popu
					outpc  chan OutpData
					mgen   int
					runs   int
					maxFit float64
				}
			*/
			popc <- gatsman.NewPop(pm, pc, ps, cross, replic, protect)
		}
	}
	close(popc)
	outpc := make(chan gatsman.OutpData, 10000)
	workerpool := make([]gatsman.Gatsworker, workers)
	for i := 0; i < workers; i++ {
		workerpool[i] = gatsman.NewGatsworker(&wg, popc, outpc, mgen, runs, maxFit)
	}
	for _, worker := range workerpool {
		go worker.Work(debug)
	}
	start := float64(time.Now().UnixNano()) / (float64(time.Second) / float64(time.Nanosecond))
	wg.Wait()
	end := float64(time.Now().UnixNano()) / (float64(time.Second) / float64(time.Nanosecond))
	fmt.Printf("************\n Duration: %.5f\n************\n", end-start)
	close(outpc)

	header := []string{"pm", "pc", "gens", fmt.Sprintf("duration: %v", end-start), fmt.Sprintf("cross: %v", cross), fmt.Sprintf("replic: %v", replic), fmt.Sprintf("%v", protect)}
	writer.Write(header)
	for {
		outp, ok := <-outpc
		if !ok {
			break
		}
		err := writer.Write(outp.ToStringArr())
		if err != nil {
			log.Fatal("cannot create file", err)
		}
	}

	defer writer.Flush()
}

func runRandom() {
	psPtr := flag.Int("ps", 100, "population size, integer")
	crossPtr := flag.Int("cross", 0, "crossover method")
	replicPtr := flag.Int("replc", 0, "replication method")
	protectPtr := flag.Bool("protect", true, "protect best bool")
	runPtr := flag.Int("runs", 100, "runs")
	maxFitPtr := flag.Float64("maxfit", 0, "max fit breakpoint")
	maxGenPtr := flag.Int("mgen", 2000, "max gens")
	inptFilePtr := flag.String("inpt", "", "inpt file")
	outpFilePtr := flag.String("outpt", fmt.Sprintf("out/outp_%v.csv", time.Now().UnixNano()), "output file")
	debugPtr := flag.Bool("debug", true, "debug flag")
	workerFlg := flag.Int("worker", 1, "number of workers to run")
	flag.Parse()

	fmt.Printf(*inptFilePtr)
	inptFile := *inptFilePtr
	gatsman.InitCMap(inptFile)
	gatsman.InitDMap()

	outpFile := *outpFilePtr
	file, err := os.Create(outpFile)
	if err != nil {
		log.Fatal("cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	ps := *psPtr
	cross := *crossPtr
	replic := *replicPtr
	protect := *protectPtr
	mgen := *maxGenPtr
	runs := *runPtr
	maxFit := *maxFitPtr
	debug := *debugPtr
	workers := *workerFlg
	var wg sync.WaitGroup
	popc := make(chan gatsman.Popu, 10000)
	popc <- gatsman.NewPop(0.03, 0.85, ps, cross, replic, protect)

	close(popc)
	outpc := make(chan gatsman.OutpData, 10000)
	workerpool := make([]gatsman.Gatsworker, workers)
	for i := 0; i < workers; i++ {
		workerpool[i] = gatsman.NewGatsworker(&wg, popc, outpc, mgen, runs, maxFit)
	}
	for _, worker := range workerpool {
		go worker.WorkRandom(debug)
	}
	start := float64(time.Now().UnixNano()) / (float64(time.Second) / float64(time.Nanosecond))
	wg.Wait()
	end := float64(time.Now().UnixNano()) / (float64(time.Second) / float64(time.Nanosecond))
	fmt.Printf("************\n Duration: %.5f\n************\n", end-start)
	close(outpc)

	header := []string{"pm", "pc", "gens", fmt.Sprintf("duration: %v", end-start), fmt.Sprintf("cross: %v", cross), fmt.Sprintf("replic: %v", replic), fmt.Sprintf("%v", protect)}
	writer.Write(header)
	for {
		outp, ok := <-outpc
		if !ok {
			break
		}
		err := writer.Write(outp.ToStringArrRand())
		if err != nil {
			log.Fatal("cannot create file", err)
		}
	}

	defer writer.Flush()
}

func createSteps(min, size float64, steps int) []float64 {
	slice := make([]float64, steps+1)
	for i := 0; i < steps; i++ {
		min += size
		slice[i] = min
	}
	return slice
}

func test1() {

	gatsman.InitCMap("05-map-10x10-36-dist42.64.txt")
	gatsman.InitDMap()
	ps := 100
	pm := 0.04
	pc := 0.85
	cross := 0
	replic := 0
	protect := false
	//pop := gatsman.NewPop(pm, pc, ps, cross, replic, protect)
	//fmt.Printf("%v\n", gatsman.GetDistance(gatsman.NewCityPair(1, 36)))

	/*fmt.Printf("%v %v\n", pop.Pop[0].Path, pop.Pop[0].Fitness)
	fmt.Printf("%v %v\n", pop.Pop[99].Path, pop.Pop[99].Fitness)*/
	var wg sync.WaitGroup
	fitsum := 0.0
	gensum := 0.0
	fitchan := make(chan float64, 100)
	genchan := make(chan int, 100)
	pathchan := make(chan []int, 100)

	runs := 10
	wg.Add(runs)
	start := float64(time.Now().UnixNano()) / (float64(time.Second) / float64(time.Nanosecond))
	for r := 0; r < runs; r++ {
		go func(wg *sync.WaitGroup, pop gatsman.Popu) {
			defer wg.Done()
			i := 0
			for i < 1000 {
				pop.Evolve()
				if pop.Pop[0].Fitness < 42.7 {
					break
				}
				i++
				/*fmt.Printf("\n-----------------------")
				for i, p := range pop.Pop {
					fmt.Printf("[%v] %v ", i, p.Fitness)
				}
				fmt.Printf("\n-----------------------")*/
			}
			fitchan <- pop.Pop[0].Fitness
			genchan <- i
			pathchan <- pop.Pop[0].Path
		}(&wg, gatsman.NewPop(pm, pc, ps, cross, replic, protect))

	}
	wg.Wait()
	end := float64(time.Now().UnixNano()) / (float64(time.Second) / float64(time.Nanosecond))
	close(fitchan)
	close(genchan)
	close(pathchan)
	for {
		fit, ok := <-fitchan
		if ok {
			fitsum += fit
		} else {
			break
		}
	}
	for {
		gen, ok := <-genchan
		if ok {
			gensum += float64(gen)
		} else {
			break
		}
	}
	fmt.Printf("last path: %v gen: %v, fit: %v time: %.4fs\n", <-pathchan, gensum/float64(runs), fitsum/float64(runs), end-start)
	/*for gen := range genchan {
		gensum += float64(gen)
	}
	for fit := range fitchan {
		fitsum += fit
	}*/
}

func parse() {
	for i := 0; i < len(os.Args); i++ {

	}
}

func initDefault() {
	gatsman.RandInitCMap(100, 50)
	gatsman.InitDMap()
}
