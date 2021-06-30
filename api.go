package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var direccionnodos []string

func nodothread(i int, chann chan string) {
	resp, err := http.Get(direccionnodos[i])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var j interface{}
	aux := json.NewDecoder(resp.Body).Decode(&j)
	if aux != nil {
		panic(aux)
	}
	fmt.Printf("%s", j)
	str := fmt.Sprintf("%v", j)
	chann <- str
}
func readnodo1(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	channel := make(chan string)
	cant := 0
	for i := 0; i < len(direccionnodos); i++ {
		go nodothread(i, channel)
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
func handleRequest() {
	http.HandleFunc("/nodo", readnodo1)
	log.Fatal(http.ListenAndServe(":80", nil))
}
func main() {
	for {
		var aux string
		fmt.Scanln(&aux)
		if aux == "go" {
			break
		} else {
			direccionnodos = append(direccionnodos, aux)
		}
	}
	handleRequest()
}
