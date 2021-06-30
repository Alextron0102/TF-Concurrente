package main

import (
	cl "TA2-GO-API/data"
	knn "TA2-GO-API/knn"
	reader "TA2-GO-API/reader"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func listarRegiones(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(cl.Regiones, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func listarProvincias(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(cl.Provincias, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func listarDistritos(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(cl.Distritos, "", " ")
	io.WriteString(res, string(jsonBytes))
}
func listarCategorias(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(cl.Categorias, "", " ")
	io.WriteString(res, string(jsonBytes))
}
func listarTipoCategorias(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(cl.Tipo_categorias, "", " ")
	io.WriteString(res, string(jsonBytes))
}
func listarSubTipoCategorias(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(cl.Sub_tipo_categorias, "", " ")
	io.WriteString(res, string(jsonBytes))
}
func knnEnpoint(res http.ResponseWriter, req *http.Request) {
	param := req.FormValue("k")
	if param == "" {
		return
	}
	k, _ := strconv.Atoi(param)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(knn.Knn(k), "", " ")
	io.WriteString(res, string(jsonBytes))
}
func handleRequest() {
	http.HandleFunc("/listar", listarRecursos)
	http.HandleFunc("/regiones", listarRegiones)
	http.HandleFunc("/provincias", listarProvincias)
	http.HandleFunc("/distritos", listarDistritos)
	http.HandleFunc("/categorias", listarCategorias)
	http.HandleFunc("/tipo_categorias", listarTipoCategorias)
	http.HandleFunc("/sub_tipo_categorias", listarSubTipoCategorias)
	http.HandleFunc("/knn", knnEnpoint)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	cl.Recursos = reader.LoadRecursos()
	knn.GetUniqueValues()
	handleRequest()
}
