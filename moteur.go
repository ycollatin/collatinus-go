/*     collatinus - moteur.go   */
package main

import (
	"fmt"
	"strings"
)

/*
  ResL :
  élément d'analyse morpho :
  lemme et morphos
*/
type Sr struct {
	Lem     *Lemme
	Morphos []string
}

// Res Collection de résultats avec plusieurs lemmes */
type Res []Sr

func restostring(r Res) string {
	var lr []string
	for _, rl := range r {
		if rl.Lem == nil {
			continue
		}
		l := fmt.Sprintf("   %s, %s : %s\n",
			strings.Join(rl.Lem.Grq, " "),
			rl.Lem.Indmorph,
			rl.Lem.Traduction)
		for _, m := range rl.Morphos {
			l = l + ("\n      " + m)
		}
		lr = append(lr, l)
	}
	return strings.Join(lr, "\n")
}

// addRes(r Res, l *Lemme, m string)
//  ajout d'une analyse lemme + morpho  à un Sr
func addRes(r Res, l *Lemme, m string) Res {
	var contient bool
	for i := 0; i < len(r); i++ {
		if r[i].Lem == l {
			contient = true
			r[i].Morphos = append(r[i].Morphos, m)
		}
	}
	if !contient {
		var rs Sr
		rs.Lem = l
		rs.Morphos = []string{m}
		r = append(r, rs)
	}
	return r
}

/*
   Lemmatisation d'une seule forme f
*/

func lemmatiseF(f string) (result Res) {
	f = strings.TrimSuffix(f, "\n")
	r := f
	d := ""

	// irréguliers
	irr := irregs[f]
	if irr != nil {
		for _, nm := range irr.lmorph {
			result = addRes(result, irr.lem, irr.grq+" "+morphos[nm])
			if irr.exclusif {
				continue
			}
		}
	}

	// romains
	if estRomain(f) {
		lin := fmt.Sprintf("%s|inv|||num.|1", f)
		romain := creeLemme(lin)
		romain.Traduction = aRomano(f)
		result = addRes(result, romain, "inv\n")
	}

	// radical-désinence
	for {
		lrad := radicaux[r]
		if len(lrad) > 0 {
			// si le radical est en -i et la dés en i-, ii peut
			//   se contracter en i
			if d > "" && d[0] == 'i' {
				nlrad := radicaux[r+"i"]
				for _, rad := range nlrad {
					lrad = append(lrad, rad)
				}
			}
			for _, rad := range lrad {
				for _, des := range rad.lemme.modele.desm[rad.num] {
					if des.gr == d {
						m := fmt.Sprintf("%s%s %s %s",
							rad.grq, des.grq, morphos[des.morpho], rad.lemme.Genre)
						result = addRes(result, rad.lemme, m)
					}
				}
			}
		}
		l := len(r) - 1
		if r == "" {
			break
		}
		d = string(r[l]) + d
		r = r[:l]
	}
	return
}

/*
   lemVars calcule toutes les variantes graphiques de f, lemmatise chacune et
   retourne avec leur analyse celles qui ont obtenu une lemmatisation
*/
func lemmatise(f string) (lsr Res, echec bool) {
	f = deramise(f)
	liste := varsF(f)
	for _, el := range liste {
		if el == "" {
			continue
		}
		lel := lemmatiseF(el)
		if len(lel) > 0 {
			for _, l := range lel {
				lsr = append(lsr, l)
			}
		} else {
			echec = true
		}
	}
	echec = len(lsr) == 0
	return
}
