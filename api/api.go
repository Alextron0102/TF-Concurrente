package main

import (
	cl "api/data"
	reader "api/reader"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var direccionnodos []string

func nodopingthread(i int, chann chan string) {
	resp, err := http.Get(direccionnodos[i] + "/ping")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var j interface{}
	aux := json.NewDecoder(resp.Body).Decode(&j)
	if aux != nil {
		panic(aux)
	}
	//fmt.Printf("%s", j)
	str := fmt.Sprintf("%v", j)
	chann <- str
}
func listarRecursos(res http.ResponseWriter, req *http.Request) {
	region := req.FormValue("region")
	provincia := req.FormValue("provincia")
	distrito := req.FormValue("distrito")
	categoria := req.FormValue("categoria")
	tipo_categoria := req.FormValue("tipo_categoria")
	sub_tipo_categoria := req.FormValue("sub_tipo_categoria")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	if region == "" && provincia == "" && distrito == "" && categoria == "" && tipo_categoria == "" && sub_tipo_categoria == "" {
		jsonBytes, _ := json.MarshalIndent(cl.Recursos, "", " ")
		io.WriteString(res, string(jsonBytes))
		return
	}
	var response []cl.Recurso
	for _, recurso := range cl.Recursos {
		if !strings.EqualFold(recurso.REGIÓN, region) && region != "" {
			continue
		}
		if !strings.EqualFold(recurso.PROVINCIA, provincia) && provincia != "" {
			continue
		}
		if !strings.EqualFold(recurso.DISTRITO, distrito) && distrito != "" {
			continue
		}
		if !strings.EqualFold(recurso.CATEGORIA, categoria) && categoria != "" {
			continue
		}
		if !strings.EqualFold(recurso.Tipo_de_Categoria, tipo_categoria) && tipo_categoria != "" {
			continue
		}
		if !strings.EqualFold(recurso.Sub_tipo_Categoria, sub_tipo_categoria) && sub_tipo_categoria != "" {
			continue
		}
		response = append(response, recurso)
	}
	jsonBytes, _ := json.MarshalIndent(response, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func readnodos(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	channel := make(chan string)
	cant := 0
	for i := 0; i < len(direccionnodos); i++ {
		go nodopingthread(i, channel)
		cant++
	}
	var responses []string
	for response := range channel {
		responses = append(responses, response)
		cant--
		if cant == 0 {
			close(channel)
		}
	}
	jsonBytes, _ := json.MarshalIndent(responses, "", " ")
	io.WriteString(res, string(jsonBytes))
}
func nodoknnthread(i int, k string, recursos []cl.Recurso, chann chan []cl.RecursoKnn) {
	jsonBytes, _ := json.MarshalIndent(recursos, "", " ")
	resp, err := http.Post(direccionnodos[i]+"/knn?k="+k, "application/json", bytes.NewBuffer(jsonBytes))
	// req, _ := http.NewRequest("POST", direccionnodos[i]+"/knn?"+k, bytes.NewBuffer(jsonBytes))
	// req.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	//resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var recursosPredict []cl.RecursoKnn
	json.NewDecoder(resp.Body).Decode(&recursosPredict)
	log.Print(recursosPredict)
	// var j interface{}
	// aux := json.NewDecoder(resp.Body).Decode(&j)
	// if aux != nil {
	// 	panic(aux)
	// }
	// fmt.Printf("%s", j)
	// str := fmt.Sprintf("%v", j)
	chann <- recursosPredict
}
func knnEnpoint(res http.ResponseWriter, req *http.Request) {
	param := req.FormValue("k")
	if param == "" {
		return
	}
	//k, _ := strconv.Atoi(param)
	//mandamos a los nodos a procesar la información:
	cantPorNodo := len(cl.Recursos) / len(direccionnodos)
	contadorRecursos := 0
	channel := make(chan []cl.RecursoKnn)
	contHilos := 0
	for i := 0; i < len(direccionnodos); i++ {
		//aqui tenemos que partir los recursos entre la cantidad de nodos:
		go nodoknnthread(i, param, cl.Recursos[contadorRecursos:contadorRecursos+cantPorNodo], channel)
		contHilos++
		contadorRecursos += cantPorNodo
	}
	var response []cl.RecursoKnn
	for recursos := range channel {
		//PrintRecurso(recurso)
		response = append(response, recursos...)
		contHilos--
		if contHilos == 0 {
			close(channel)
		}
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(response, "", " ")
	io.WriteString(res, string(jsonBytes))
}
func handleRequest() {
	http.HandleFunc("/listar", listarRecursos)
	http.HandleFunc("/knn", knnEnpoint)
	http.HandleFunc("/pingnodos", readnodos)
	log.Fatal(http.ListenAndServe(":80", nil))
}
func main() {
	log.Println("inicializando...")
	cl.Recursos = reader.LoadRecursos()
	log.Println("escuchando...")
	fmt.Println("Ingrese las direcciones de nodos (ingrese go para parar):")
	for {
		var aux string
		fmt.Scanln(&aux)
		if aux == "go" {
			break
		} else {
			direccionnodos = append(direccionnodos, aux)
		}
	}
	log.Println("Nodos:")
	for _, dir := range direccionnodos {
		log.Println(dir)
	}
	handleRequest()
}
