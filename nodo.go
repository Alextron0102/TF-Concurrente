package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func hola(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent("hola desde un nodo4", "", " ")
	io.WriteString(res, string(jsonBytes))
}
func handleRequest() {
	http.HandleFunc("/hola", hola)
	log.Fatal(http.ListenAndServe(":84", nil))
}
func main() {
	handleRequest()
}
