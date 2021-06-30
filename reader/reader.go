package reader

//DATASET:
//https://www.datosabiertos.gob.pe/dataset/inventario-de-recursos-tur%C3%ADsticos

import (
	//"bufio"
	cl "TA2-GO-API/data"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	r.Comma = ','
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}
	data, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return data, nil
}
func convertThread(tuples [][]string, ch chan cl.Recurso) {
	var latitud, longitud *float64 = new(float64), new(float64)
	for _, value := range tuples {
		codigo, err := strconv.Atoi(strings.Replace(strings.TrimSpace(value[3]), ",", "", -1))
		check(err)
		if value[9] != "" {
			aux, err := strconv.ParseFloat(strings.TrimSpace(value[9]), 64)
			check(err)
			latitud = &aux
		} else {
			latitud = nil
		}
		if value[10] != "" {
			aux, err := strconv.ParseFloat(strings.TrimSpace(value[10]), 64)
			check(err)
			longitud = &aux
		} else {
			longitud = nil
		}
		recursosingle := cl.Recurso{
			REGIÓN:             value[0],
			PROVINCIA:          value[1],
			DISTRITO:           value[2],
			Codigo_del_Recurso: codigo,
			Nombre_del_Recurso: value[4],
			CATEGORIA:          value[5],
			Tipo_de_Categoria:  value[6],
			Sub_tipo_Categoria: value[7],
			URL:                value[8],
			LATITUD:            latitud,
			LONGITUD:           longitud,
		}
		//put to channel
		ch <- recursosingle
		//recursos = append(recursos, recursosingle)
	}
	close(ch)
}
func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func LoadRecursos() []cl.Recurso {
	data, err := readCSVFromUrl("https://raw.githubusercontent.com/Alextron0102/TA2-Go-API/main/files/Inventario_recursos_turisticos.csv")
	check(err)
	channels := make([]chan cl.Recurso, cl.NUM_CPU+1)
	limit := len(data) / cl.NUM_CPU
	fmt.Print("lineas en total: ")
	fmt.Println(len(data))
	iteratoraux := 0
	for i := 0; i < len(data); i += limit {
		chunk := data[i:Min(i+limit, len(data))]
		channels[iteratoraux] = make(chan cl.Recurso)
		go convertThread(chunk, channels[iteratoraux])
		iteratoraux++
	}
	var recursos []cl.Recurso
	for _, channel := range channels {
		for recurso := range channel {
			//PrintRecurso(recurso)
			recursos = append(recursos, recurso)
		}
	}
	return recursos
}
