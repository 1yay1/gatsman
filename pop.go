package gatsman

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const greedyCross int = 0

const double50 int = 0
const tournament int = 1
const topHalf int = 2

//Popu struct
type Popu struct {
	Pm      float64
	Pc      float64
	cross   int
	replic  int
	protect bool

	ps  int
	Pop []*Tsman
}

//String method for Popu
func (p Popu) String() string {
	return fmt.Sprintf("pm: %v, pc: %v, cross: %v, replic: %v, protect :%v", p.Pm, p.Pc, p.cross, p.replic, p.protect)
}

//NewPop creates a new Population element
func NewPop(pm, pc float64, ps, cross, replic int, protect bool) Popu {
	popArr := make([]*Tsman, ps)
	for i := 0; i < ps; i++ {
		popArr[i] = NewTsman()
	}
	return Popu{pm, pc, cross, replic, protect, ps, popArr}
}

//ByFit implements sort interface
type ByFit []*Tsman

func (a ByFit) Len() int           { return len(a) }
func (a ByFit) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFit) Less(i, j int) bool { return a[i].Fitness < a[j].Fitness }

//NewCopy create a new population with the same params as the og
func (p *Popu) NewCopy() Popu {
	return NewPop(p.Pm, p.Pc, p.ps, p.cross, p.replic, p.protect)
}

//Evolve ...CityPair
func (p *Popu) Evolve() {
	rand.Seed(time.Now().UnixNano())
	var best Tsman
	var newPopu = p.Pop
	var cmax = int(p.Pc * float64(p.ps))
	var mmax = int(p.Pm * float64(p.ps) * float64(CityCount()))
	if p.protect {
		best = *p.Pop[0]
	}

	for i := 0; i < cmax; i++ {
		iarr := rand.Perm(p.ps)
		if p.cross == greedyCross {
			newPopu[i] = GreedyCross(p.Pop[iarr[0]], p.Pop[iarr[1]])
		} else {
			//possible different cross func here!
			newPopu[i] = GreddyQuadCross(p.Pop[iarr[0]], p.Pop[iarr[1]])
		}
	}
	for i := cmax; i < p.ps; i++ {
		newPopu[i] = p.Pop[rand.Intn(p.ps)]
		newPopu[i].Fitness = CalcFit(newPopu[i].Path)
	}

	for m := 0; m < mmax; m++ {
		var idx int
		if p.protect {
			idx = 1 + rand.Intn(p.ps-1)
		} else {
			idx = rand.Intn(p.ps)
		}
		newPopu[idx].Mutate()
		newPopu[idx].Fitness = CalcFit(p.Pop[idx].Path)
	}

	sort.Sort(ByFit(newPopu))
	if p.protect {
		newPopu[p.ps-1] = &best
		sort.Sort(ByFit(newPopu))
	}

	if p.replic == double50 {
		first := *newPopu[0]
		next := *newPopu[1]
		idxs := rand.Perm(p.ps)
		for i := p.ps / 2; i < p.ps; i++ {
			newPopu[i] = newPopu[idxs[i]]
		}
		i := 0
		for ; i < p.ps/4; i++ {
			newPopu[i] = &first
		}
		for ; i < p.ps/2; i++ {
			newPopu[i] = &next
		}

	} else if p.replic == tournament {
		for i := p.ps / 2; i < p.ps; i++ {
			ts1 := newPopu[rand.Intn(p.ps)]
			ts2 := newPopu[rand.Intn(p.ps)]
			ts3 := newPopu[rand.Intn(p.ps)]

			if ts1.Fitness < ts2.Fitness {
				if ts1.Fitness < ts3.Fitness {
					newPopu[i] = ts1
				} else {
					newPopu[i] = ts3
				}
			} else {
				if ts2.Fitness < ts3.Fitness {
					newPopu[i] = ts2
				} else {
					newPopu[i] = ts3
				}
			}
		}
	} else if p.replic == topHalf {
		//double top 50%
		for i := 0; i < p.ps/2; i++ {
			newPopu[p.ps/2+i] = newPopu[i]
		}
	} else {
		for i := p.ps - 1; i >= 0; i -= p.ps / 10 {
			for j := 0; j < 10; j++ {
				newPopu[i-j] = newPopu[j]
			}
		}
	}
	sort.Sort(ByFit(newPopu))
	p.Pop = newPopu
}
