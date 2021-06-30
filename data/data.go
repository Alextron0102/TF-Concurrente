package data

import (
	"runtime"
)

var Recursos []Recurso
var Regiones, Provincias, Distritos, Categorias, Tipo_categorias, Sub_tipo_categorias []string

type Recurso struct {
	REGIÓN             string
	PROVINCIA          string
	DISTRITO           string
	Codigo_del_Recurso int
	Nombre_del_Recurso string
	CATEGORIA          string
	Tipo_de_Categoria  string
	Sub_tipo_Categoria string
	URL                string
	LATITUD            *float64
	LONGITUD           *float64
}

type RecursoKnn struct {
	Recurso                    *Recurso
	REGIÓN_predict             string
	PROVINCIA_predict          string
	DISTRITO_predict           string
	CATEGORIA_predict          string
	Tipo_de_Categoria_predict  string
	Sub_tipo_Categoria_predict string
}

var NUM_CPU = runtime.NumCPU()-1
