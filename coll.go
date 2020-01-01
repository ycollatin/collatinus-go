//   collatinus  - lemmatisation et analyse morpho de textes latins
package main

/*
	XXX - bogues

	TODO
	- écrire un lisdata() pour que les données du module soient toutes lues.
	  Seuls les lemmes sont lus. Il faudrait ajouter
	   . irréguliers
	   . vargraph
	   . modeles
	- moteur.go : il ne devrait pas y avoir de désinence nil dans la
	  boucle lemmatise()
	- paramètres supplémentaires possibles :
        -w sortie html
		-t traduction
		-l tri alpha des lemmes
		-m tri alpha des mots
		-e rejet des échecs en fin
		-c calcul des stats (fréquence des formes, des lemmes, des morphos, des pos)
		-o <fichier de sortie>
		-s serveur sur le port indiqué, ou un port par défaut
	- modeles.la : nombreuses désinences héritées redéfinies à l'identique
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var dat bool

const version = "Alpha"

var modules []string

// lecture des données et affichage des effectifs
func data() {
	if dat {
		return
	}
	lismorphos()
	fmt.Println(len(morphos), "morphos")
	lismodeles()
	fmt.Println(len(modeles), "modèles")
	var nc string
	if len(module) > 0 {
		fmt.Println("module", module)
		nc = "data/" + module + "/"
		lisLemmes(nc + "lemmes.la")
		lisTraductions(nc + "lemmes.fr")
		//lisIrregs(nc + "irregs.la")
		lisExp(nc + "vargraph.la")
	}
	lisExp("data/vargraph.la")
	lisLemmes("data/lemmes.la")
	lisTraductions("data/lemmes.fr")
	fmt.Println(len(lemmes), "lemmes")
	ajRadicaux()
	fmt.Println(len(radicaux), "radicaux")
	if len(module) > 0 {
		lisIrregs("data/" + module + "/irregs.la")
	}
	lisIrregs("data/irregs.la")
	fmt.Println(len(irregs), "irréguliers")
	fmt.Println(len(lexp), "variantes graphiques\n")
	dat = true
}

func interact() {
	data()
	lecteur := bufio.NewReader(os.Stdin)
	var f string
	for f != "x" {
		fmt.Print("> ")
		f, _ = lecteur.ReadString('\n')
		lm := mots(f)
		for _, m := range lm {
			m = strings.Trim(m, "\n\r")
			if m == "x" {
				return
			}
			ret, echec := lemmatise(m)
			if echec {
				ret, echec = lemmatise(majminmaj(f))
			}
			for _, r := range ret {
				fmt.Printf("   %s, %s, %s : %s\n",
					strings.Join(r.Lem.Grq, " "),
					r.Lem.Pos,
					r.Lem.Indmorph,
					r.Lem.Traduction)
				for _, m := range r.Morphos {
					fmt.Println("      ", m)
				}
			}
		}
	}
}

var module string

func main() {

	var h bool
	flag.BoolVar(&h, "h", false, "usage")
	var ar = flag.Bool("a", false, "serveur tcp. avec l'option -p")
	var nf = flag.String("f", "", "nom du fichier à lemmatiser")
	var nfs = flag.String("o", "", "nom du fichier de résultats")
	var m = flag.String("m", "", "module lexical et morphologique")
	var s = flag.Bool("s", false, "Collatinus est serveur, port par défaut : 8080")
	var p = flag.String("p", "8080", "Port du ѕerveur. doit être précédéz de l'option -s")
	var i = flag.Bool("i", false, "mode interactif")

	module = *m

	flag.Parse()
	fmt.Println("Collatinus Go")
	fmt.Println("© Yves Ouvrard, GPL3")
	if *m > "" {
		fmt.Println("module", *m)
	}

	if h {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *i {
		interact()
		os.Exit(0)
	}

	if *nf != "" {
		data()
		// analyse du contenu du fichier en argument
		ll := lignes(*nf)
		var sortie *os.File
		sf := *nfs != ""
		if sf {
			sortie, _ = os.Create(*nfs)
		}
		for _, l := range ll {
			lm := mots(l)
			for _, m := range lm {
				an, echec := lemmatise(m)
				if echec {
					an, echec = lemmatise(majminmaj(m))
				}
				if echec {
					if sf {
						sortie.WriteString("\n" + m + " échec")
					} else {
						fmt.Println(m, "échec")
					}
				} else if sf {
					sortie.WriteString("\n" + restostring(an))
				} else {
					fmt.Print(restostring(an))
				}
			}
		}
		if sf {
			sortie.Close()
		} else {
			fmt.Println("\n----oOo----")
		}
		os.Exit(0)
	}
	if *ar {
		data()
		arbos(*p)
		os.Exit(0)
	}
	if *s {
		data()
		serveur(*p)
	} else {
		interact()
	}
}
