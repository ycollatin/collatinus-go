package main

import (
	"fmt"
	//"log"
	"strings"
)

type Lemme struct {
	Grq,
	Gr []string
	cle,
	Indmorph,
	Pos,
	Genre,
	Traduction string
	nh,
	Freq int // fréquence en dernier champ
	radicaux map[int][]*Rad
	modele   *Modele
	// TODO ajouter une map de traductions
}

func (l Lemme) doc() string {
	var lr []string
	for nr, rr := range l.radicaux {
		for _, r := range rr {
			lr = append(lr, fmt.Sprintf("n° %d %s", nr, r.grq))
		}
	}
	return fmt.Sprintf("clé %s %s, %s-%s : %s\n  rad. %s",
		l.cle,
		strings.Join(l.Grq, ","), l.Indmorph, l.Pos, l.Traduction,
		strings.Join(lr, "\n       "))
}

func (l Lemme) habetRad(r string, n int) bool {
	for _, rad := range l.radicaux[n] {
		if rad.gr == r {
			return true
		}
	}
	return false
}

var lemmes = make(map[string]*Lemme)

func recupNh(k string) (string, string) {
	der := string(k[len(k)-1])
	if der == "2" || der == "3" || der == "4" {
		k = strings.TrimSuffix(k, der)
	}
	return k, der
}

// creeLemme(l string)
// créateur de lemme à partir de la ligne l de lemmes.la
func creeLemme(l string) *Lemme {
	var lem *Lemme = new(Lemme)
	lem.radicaux = make(map[int][]*Rad)
	eclats := strings.Split(l, "|")
	for i, e := range eclats {
		switch i {
		case 0:  // forme canonique
			eclg := strings.Split(e, "=")
			// supprimer et affecter le numéro d'homonymie
			var lgrq, lnh string
			lgrq, lnh = recupNh(eclg[0])
			lem.nh = strtoint(lnh)
			lem.Grq = append(lem.Grq, lgrq)
			lem.Gr = append(lem.Gr, atone(eclg[0]))
			if len(eclg) > 1 {
				eclk := strings.Split(eclg[1], ",")
				for i, eclki := range eclk {
					lem.Grq = append(lem.Grq, eclki)
					lem.Gr = append(lem.Gr, atone(lem.Grq[i]))
				}
			}
			lem.cle = atone(lem.Grq[0])
			if lem.nh > 1 {
				lem.cle = lem.cle + lnh
			}

		case 1:  // modèle
			lem.modele = modeles[e]
			if strings.ToLower(lem.cle) != lem.cle {
				lem.Pos = "NP"
			} else {
				lem.Pos = lem.modele.pos
			}

		case 2, 3: // 2 radicaux gén/pf, 3 supin
			if e > "" {
				eclR := strings.Split(e, ",")
				for _, r := range eclR {
					rgr := deramAtone(r)
					rad := new(Rad)
					rad.grq = r
					rad.gr = rgr
					rad.num = i - 1
					rad.lemme = lem
					lem.radicaux[rad.num] = append(lem.radicaux[rad.num], rad)
				}
			}
			// radicaux calculés
			for _, grad := range lem.modele.lgenR {
				// si un radical de même num a été donné
				// on ne tient pas compte du calculé
				if len(lem.radicaux[grad.num]) > 0 {
					continue
				}
				for _, grq := range lem.Grq {
					rgrq := grad.rad(grq)
					rgr := deramAtone(rgrq)
					if !lem.habetRad(rgr, grad.num) {
						rad := new(Rad)
						rad.grq = rgrq
						rad.gr = rgr
						rad.num = grad.num
						rad.lemme = lem
						lem.radicaux[rad.num] = append(lem.radicaux[rad.num], rad)
					}
				}
			}
		case 4: // indMorph
			lem.Indmorph = e
			if strings.HasSuffix(lem.Indmorph, " f.") {
				lem.Genre = "féminin"
			} else if strings.HasSuffix(lem.Indmorph, " m.") {
				lem.Genre = "masculin"
			} else if strings.HasSuffix(lem.Indmorph, " n.") {
				lem.Genre = "neutre"
			}
			// pos des prépositions, négations et adverbes
			cacc := contient(lem.Indmorph, "+ acc.")
			cabl := contient(lem.Indmorph, "+ abl.")
			switch {
			case cacc && cabl:
				lem.Pos = "prepAA"
			case cabl:
				lem.Pos = "prepAbl"
			case cacc:
				lem.Pos = "prepAcc"
			case contient(lem.Indmorph, "neg."):
				lem.Pos = "neg"
			case contient(lem.Indmorph, "adv"):
				lem.Pos = "Adv"
			case contient(lem.Indmorph, "interj."):
				lem.Pos = "intj"
			}
		case 5: // fréquence
			lem.Freq = strtoint(e)
		}
	}
	return lem
}

// lisLemmes()
// lecteur de la base de lemmes
func lisLemmes(nf string) {
	ll := Lignes(nf)
	for _, l := range ll {
		lem := creeLemme(l)
		// si le lemme existe déjà, passer)
		present, _ := lemmes[lem.cle]
		if present == nil {
			lemmes[lem.cle] = lem
		}
	}
}

// lisTraductions()
// lis les traductions à martir d'un fichier lemmes.??
func lisTraductions(nf string) {
	ll := Lignes(nf)
	for _, l := range ll {
		ecl := strings.Split(l, ":")
		cle := ecl[0]
		lem := lemmes[cle]
		if lem != nil && lem.Traduction == "" {
			lem.Traduction = strings.Join(ecl[1:], ":")
		} /* else { log.Println(l) } */
	}
}
