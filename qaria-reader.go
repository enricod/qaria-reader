package main

/**
 * programma a linea di comando per lettura da pagine html dei dati di inquinamento
 *
 */
import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"github.com/enricod/qaria-model"
)


type DbConf struct {
	User     string
	Password string
}

func main() {

	fmt.Printf("Lettore dati inquinamento per lombardia\n")
	dbusername := flag.String("dbusername", "root", "Utente per connessione a db")
	flag.Parse()

	dbconf := DbConf{*dbusername, "root"}
	stazioni := qariamodel.ElencoStazioni()

	for _, s := range stazioni {
		misure, err := LeggiMisure(s)
		if err != nil {
			fmt.Printf("errore %v", err)
		} else {
			fmt.Printf("misure lette %v\n", misure)
			salvaInDb(dbconf, misure)
		}
	}
}



func salvaInDb(dbconf DbConf, misure []qariamodel.Misura) {
	if db, err := sql.Open("mysql",
		dbconf.User+":"+
			dbconf.Password+
			"@tcp(127.0.0.1:3306)/qaria"); err != nil {
		log.Fatal(err)
	} else {
		if stmt, err2 := db.Prepare("INSERT INTO misura(inquinante, valore, stazioneId, dataStr) VALUES(?, ?, ?, ?)"); err2 != nil {
			log.Fatal(err2)
		} else {
			for _, m := range misure {
				if res, err3 := stmt.Exec(m.Inquinante, m.Valore, m.StazioneId, m.DataMisura); err3 != nil {
					rowCnt, err4 := res.RowsAffected()
					if err4 != nil {
						log.Fatal(err)
					} else {
						log.Printf("dati inseriti %v\n", rowCnt)
					}
				}
			}
		}
		defer db.Close()
	}
}

func LeggiMisure(s qariamodel.Stazione) ([]qariamodel.Misura, error) {

	fmt.Printf("STAZIONE %v, \t URL=%v\n", s.Nome, s.Url)
	if resp, err := http.Get(s.Url); err == nil {
		if htmlData, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			bodyStr := string(htmlData)

			result, _ := EstraiMisure(s, bodyStr)
			return result, nil
		} else {
			fmt.Printf("%v", err2)
			return nil, err2
		}
	} else {
		fmt.Printf("%v", err)
		return nil, err
	}

}

func EstraiMisure(s qariamodel.Stazione, htmlStr string) ([]qariamodel.Misura, error) {
	inquinanti := strings.Split(s.Inquinanti, ",")
	dataMisura, _ := estraiDataDaHTML(htmlStr)
	// fmt.Printf("\t htmlStr = %v\n", htmlStr)
	var result []qariamodel.Misura
	for _, inq := range inquinanti {
		// fmt.Printf("\t inq = %v\n", inq)

		i := strings.Index(htmlStr, "> " + inq + "<")
		// fmt.Printf("\t i = %v\n", i)
		if i > 0 {
			s2 := htmlStr[i:len(htmlStr)]
			r, _ := regexp.Compile("([0-9.]+)&nbsp;&nbsp; <")
			// fmt.Printf("str = %v", r.FindStringSubmatch(s2)[1])
			val, _ := strconv.ParseFloat(r.FindStringSubmatch(s2)[1], 32)

			misura := qariamodel.Misura{
				StazioneId: s.StazioneId,
				Inquinante: inq,
				Valore:     val,
				DataMisura: dataMisura,
			}

			result = append(result, misura)

		} else {

		}
	}
	return result, nil
}

// nel file html, la data ha la forma 1 dicembre 2016.
// dobbiamo convertirla in 20161201
func convertiData(dataHTML string) string {
	pezzi := strings.Split(dataHTML, " ")
	mese := "01"
	switch pezzi[1] {
	case "gennaio":
		mese = "01"
	case "febbraio":
		mese = "02"
	case "marzo":
		mese = "03"
	case "aprile":
		mese = "04"
	case "maggio":
		mese = "05"
	case "giugno":
		mese = "06"
	case "luglio":
		mese = "07"
	case "agosto":
		mese = "08"
	case "settembre":
		mese = "09"
	case "ottobre":
		mese = "10"
	case "novembre":
		mese = "11"
	case "dicembre":
		mese = "12"
	}

	if len(pezzi[0]) < 2 {
		return pezzi[2] + mese + "0" + pezzi[0]
	}
	return pezzi[2] + mese + pezzi[0]
}

func estraiDataDaHTML(htmlStr string) (string, error) {
	r, _ := regexp.Compile("<span style=\"font-size:20pt;\">Gli inquinanti monitorati <b> ([a-zA-Z0-9 ]+) </b></span>")
	val := r.FindStringSubmatch(htmlStr)[1]
	return convertiData(val), nil
}
