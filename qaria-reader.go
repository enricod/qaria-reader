package main

import (
	"fmt"
	"log"
	"strings"
	"regexp"
	"strconv"
	"net/http"
	"io/ioutil"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	)

type Stazione struct {
  StazioneId int
  Nome string
  Url string
  Inquinanti string
}

type Misura struct {
    DataMisura string
    Inquinante string
    StazioneId int
    Valore float64
}


func main() {

	fmt.Printf("Lettore dati inquinamento per lombardia\n")
	s1 := Stazione{ StazioneId:661, 
			Nome:"Rezzato",
			Inquinanti:"PM10,NO2,CO",
			Url:"http://www2.arpalombardia.it/sites/QAria/_layouts/15/QAria/DettaglioStazione.aspx?IdStaz=661"}

	stazioni := []Stazione{ s1 }


	misure, err := LeggiMisure(stazioni[0])
	if (err != nil) {
		fmt.Printf("errore %v",  err)
	} else {
		fmt.Printf("misure lette %v\n",  misure)
		salvaInDb( misure )
	}
}


func salvaInDb(misure []Misura) {
	dbUsername := "root"
	if db, err := sql.Open("mysql", dbUsername + ":root@tcp(127.0.0.1:3306)/qaria"); err != nil {
		log.Fatal(err)
	} else {
		if stmt, err2 := db.Prepare("INSERT INTO misura(inquinante, valore, stazioneId, dataStr) VALUES(?, ?, ?, ?)"); err2 != nil {
			
			log.Fatal(err2)
		} else {
			for _, m := range misure {
				if res, err3 := stmt.Exec(m.Inquinante, m.Valore, m.StazioneId, m.DataMisura); 
							err3 != nil {
					rowCnt, err4 := res.RowsAffected()
					if err4 != nil {
						log.Fatal(err)
					} else {
						log.Printf("dati inseriti %v\n", rowCnt)
					}
				}
			}
			
		}
	}
	
}

func LeggiMisure( s Stazione) ([]Misura, error) {
	fmt.Printf("STAZIONE nome=%v,\n\t URL=%v\n", s.Nome, s.Url)
	if resp, err :=  http.Get(s.Url); err == nil {
		if htmlData, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			bodyStr := string(htmlData)
			
			result, _ := EstraiMisure(s,bodyStr)
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




func EstraiMisure(s Stazione, htmlStr string) ([]Misura, error) {
	inquinanti := strings.Split(s.Inquinanti, ",")
	dataMisura, _ := estraiDataDaHTML(htmlStr)
	var result []Misura
	for _, inq := range inquinanti {
		i := strings.Index(htmlStr, "> "+inq+"  <")
		if i > 0 {
			
			s2 := htmlStr[i:len(htmlStr)]
			r, _ := regexp.Compile("([0-9.]+)&nbsp;&nbsp; <")
			val, _ := strconv.ParseFloat(r.FindStringSubmatch(s2)[1], 32)

			misura := Misura{
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