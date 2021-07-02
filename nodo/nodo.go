package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	cl "nodo/data"
	knn "nodo/knn"
	reader "nodo/reader"
	"os"
	"strconv"
)

var port string

func knnEnpoint(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		kparam := req.FormValue("k")
		if kparam == "" {
			return
		}
		k, _ := strconv.Atoi(kparam)
		var recursos []cl.Recurso
		json.NewDecoder(req.Body).Decode(&recursos)
		//log.Println(recursos)
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		jsonBytes, _ := json.MarshalIndent(knn.Knn(k, recursos), "", " ")
		io.WriteString(res, string(jsonBytes))
	}
}

//funcion para ver el puerto en el que está corriendo:
func hola(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent("hola desde el puerto "+port, "", " ")
	io.WriteString(res, string(jsonBytes))
}
func handleRequest() {
	http.HandleFunc("/ping", hola)
	http.HandleFunc("/knn", knnEnpoint)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func main() {
	log.Println("inicializando...")
	cl.Recursos = reader.LoadRecursos()
	knn.GetUniqueValues()
	log.Println("escuchando...")
	//port = "80"
	port = os.Getenv("PORT")
	handleRequest()
}
