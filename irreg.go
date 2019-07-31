package main

import (
	"strings"
)

type Irr struct {
	grq,
	gr string
	lem    *Lemme
	lmorph []int
}

func creeIrr(l string) (irr *Irr) {
	ecl := strings.Split(l, ":")
	irr = new(Irr)
	irr.grq = ecl[0]
	irr.gr = atone(irr.grq)
	irr.lem = lemmes[ecl[1]]
	irr.lmorph = listei(ecl[2])
	return irr
}

var irregs = make(map[string]*Irr)

func lisIrregs(nf string) {
	//ll := lignes("data/irregs.la")
	ll := lignes(nf)
	for _, l := range ll {
		nirr := creeIrr(l)
		if irregs[nirr.gr] == nil {
			irregs[nirr.gr] = nirr
		}
	}
}

/*
   lem := lemmes[cle]
   if lem != nil && lem.Traduction == "" {
       lem.Traduction = strings.Join(ecl[1:], ":")
   }
*/
