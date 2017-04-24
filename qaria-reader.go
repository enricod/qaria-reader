package main

import (
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"net/http"
	"io/ioutil"
	)


func main() {
	fmt.Printf("lettore dati inquinamento per lombardia\n")
	s1 := Stazione{ StazioneId:661, 
		Nome:"Rezzato",
		Inquinanti:"PM10,NO2,CO",
		Url:"http://www2.arpalombardia.it/sites/QAria/_layouts/15/QAria/DettaglioStazione.aspx?IdStaz=661"}

	misure, err := LeggiMisure(s1)
	if (err != nil) {
		fmt.Printf("errore %v",  err)
	} else {
		fmt.Printf("misure lette %v\n",  misure)
	}
}

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