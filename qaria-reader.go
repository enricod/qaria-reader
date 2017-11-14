package main

/**
 * programma a linea di comando per lettura da pagine html dei dati di inquinamento
 *
 */
import (
	//"database/sql"
	"flag"
	"fmt"
	"github.com/enricod/qaria-model"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	//"log"
	"net/http"
	"strconv"
	"time"
)

type DbConf struct {
	User     string
	Password string
}

func main() {
	// directory dove salvare il file html
	outputdir := flag.String("d", ".", "outputdir")
	flag.Parse()

	fmt.Printf("Lettore dati inquinamento per lombardia\n")
	fmt.Println(fmt.Sprintf("\t cartella dove salverò i dati: %v", *outputdir))

	//dbconf := DbConf{*dbusername, "root"}
	stazioni := qariamodel.ElencoStazioni()

	for _, s := range stazioni {
		filename, err := LeggiPaginaWeb(*outputdir, s)
		if err != nil {
			fmt.Printf("errore %v", err)
		} else {
			fmt.Printf("creato file  %v\n", filename)
		}
	}
}

//func salvaInDb(dbconf DbConf, misure []qariamodel.Misura) {
//	if db, err := sql.Open("mysql",
//		dbconf.User+":"+
//			dbconf.Password+
//			"@tcp(127.0.0.1:3306)/qaria"); err != nil {
//		log.Fatal(err)
//	} else {
//		if stmt, err2 := db.Prepare("INSERT INTO misura(inquinante, valore, stazioneId, dataStr) VALUES(?, ?, ?, ?)"); err2 != nil {
//			log.Fatal(err2)
//		} else {
//			for _, m := range misure {
//				if res, err3 := stmt.Exec(m.Inquinante, m.Valore, m.StazioneId, m.DataMisura); err3 != nil {
//					rowCnt, err4 := res.RowsAffected()
//					if err4 != nil {
//						log.Fatal(err)
//					} else {
//						log.Printf("dati inseriti %v\n", rowCnt)
//					}
//				}
//			}
//		}
//		defer db.Close()
//	}
//}

/* return il nome del file in cui è stata salvata la pagina web
 */
func LeggiPaginaWeb(dir string, s qariamodel.Stazione) (string, error) {
	t := time.Now()
	fmt.Println()

	// nome del file dove salvare la pagina HTML
	filename := dir + "/" + t.Format("20160102150405") + "_" + strconv.Itoa(s.StazioneId) + `.html`

	fmt.Printf("lettura stazione %v, URL=%v\n", s.Nome, s.Url)
	if resp, err := http.Get(s.Url); err == nil {
		if htmlData, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			bodyStr := string(htmlData)
			err = ioutil.WriteFile(filename, []byte(bodyStr), 0644)
		}
	}
	return filename, nil
}
