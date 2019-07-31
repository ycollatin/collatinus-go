/*       collatinus - arbos.go  */

package main

import (
	"bufio"
	//"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

/* types pour l'envoi des résultats de lemmatisation  */

/* une forme */
type Asr struct {
	Lem     string
	Freq    int
	Pos     string
	Intrans bool
	Morphos []string
}

/*   un mot */
type Sm struct {
	Gr   string
	Llem []Asr
}

/*  un texte */
type St []Sm

func lemarbos(t string) (st St) {
	lm := mots(t)
	for _, m := range lm {
		var sm Sm
		sm.Gr = m
		an, echec := lemmatise(m)
		if echec {
			an, _ = lemmatise(majminmaj(m))
		}
		// conversion Res -> Sm
		for _, r := range an {
			var nsm Asr
			nsm.Lem = r.Lem.Gr[0]
			nsm.Freq = r.Lem.Freq
			nsm.Pos = r.Lem.Pos
			nsm.Intrans = contient(r.Lem.Indmorph, "intr.")
			nsm.Morphos = r.Morphos
			sm.Llem = append(sm.Llem, nsm)
		}
		// conversion Sm -> St
		st = append(st, sm)
	}
	return
}

func arbos(port string) {
	fmt.Println("écoute port", port)
	ln, _ := net.Listen("tcp", ":"+port)
	for {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n')
		if len(message) < 1 {
			continue
		}
		fmt.Println(string(message))
		resultat := lemarbos(message)
		encoder := gob.NewEncoder(conn)
		encoder.Encode(resultat)
	}
	ln.Close()
}
