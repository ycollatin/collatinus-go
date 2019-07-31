package main

import (
	"fmt"
	"strconv"
	"strings"
)

var lenm int
var morphos = make(map[int]string)

func lismorphos() {
	ll := lignes("data/morphos.fr")
	for i := 0; i < len(ll); i++ {
		ecl := strings.Split(ll[i], ":")
		if len(ecl) > 1 {
			n, _ := strconv.Atoi(ecl[0])
			morphos[n] = ecl[1]
		}
	}
	lenm = len(morphos)
}

func morpho(i int) string {
	if i > lenm {
		return fmt.Sprintf("err. il n'y a que %d morphos\n", lenm)
	}
	return morphos[i]
}
