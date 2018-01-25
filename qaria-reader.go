package main

/**
 * programma a linea di comando per lettura da pagine html dei dati di inquinamento
 *
 */
import (
	//"database/sql"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/enricod/qaria-model"
	_ "github.com/go-sql-driver/mysql"
	//"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	// directory dove salvare il file html
	outputdir := flag.String("d", ".", "outputdir")
	flag.Parse()

	fmt.Printf("Lettore dati inquinamento per lombardia\n")
	fmt.Println(fmt.Sprintf("\t cartella dove salverò i dati: %v", *outputdir))

	stazioni := qariamodel.ElencoStazioni()

	for _, s := range stazioni {
		filename, err := leggiPaginaWeb(*outputdir, s)
		if err != nil {
			fmt.Printf("errore %v", err)
		} else {
			fmt.Printf("creato file  %v\n", filename)
		}
	}
}

// leggiPaginaWeb return il nome del file in cui è stata salvata la pagina web
func leggiPaginaWeb(dir string, s qariamodel.Stazione) (string, error) {
	t := time.Now()
	fmt.Println()

	// nome del file dove salvare la pagina HTML
	filename := dir + "/" + t.Format(time.RFC3339) + "_" + strconv.Itoa(s.StazioneID) + `.html`

	fmt.Printf("lettura stazione %v, URL=%v\n", s.Nome, s.URL)
	if resp, err := http.Get(s.URL); err == nil {
		if htmlData, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			bodyStr := string(htmlData)
			err = ioutil.WriteFile(filename, []byte(bodyStr), 0644)
		}
	}
	return filename, nil
}
