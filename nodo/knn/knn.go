package knn

import (
	"math"
	cl "nodo/data"
)

type DistanciaSt struct {
	origen    *cl.Recurso
	destino   *cl.Recurso
	distancia float64
}

func distancia(rec1a, rec2a cl.Recurso) DistanciaSt {
	return DistanciaSt{
		origen:    &rec1a,
		destino:   &rec2a,
		distancia: math.Sqrt(math.Pow(*rec1a.LATITUD-*rec2a.LATITUD, 2) + math.Pow(*rec1a.LONGITUD-*rec2a.LONGITUD, 2)),
	}
}

func getindicemayor(arr []int) int {
	indice := 0
	mayor := arr[indice]
	for i := 1; i < len(arr); i++ {
		if arr[i] > mayor {
			mayor = arr[i]
			indice = i
		}
	}
	return indice
}
func searchindex(array []string, str string) int {
	for i, ele := range array {
		if str == ele {
			return i
		}
	}
	return -1
}
func predictRecurso(origen cl.Recurso, k int, recursoKnn chan cl.RecursoKnn) {
	var distancias []DistanciaSt
	for _, recurso := range cl.Recursos {
		if recurso.LATITUD != nil && recurso.LONGITUD != nil {
			distancias = append(distancias, distancia(origen, recurso))
		}
		//Si el recurso no tiene latitud ni longitud no podemos aplicar el algoritmo por distancia
	}
	for i := 0; i < len(distancias); i++ {
		for j := i; j < len(distancias); j++ {
			if distancias[i].distancia > distancias[j].distancia {
				aux := distancias[i]
				distancias[i] = distancias[j]
				distancias[j] = aux
			}
		}
	}
	//predecimos el dato a base de los k primeros recursos de la lista:
	frecuencias := make([][]int, 6)
	frecuencias[0] = make([]int, len(cl.Regiones))
	frecuencias[1] = make([]int, len(cl.Provincias))
	frecuencias[2] = make([]int, len(cl.Distritos))
	frecuencias[3] = make([]int, len(cl.Categorias))
	frecuencias[4] = make([]int, len(cl.Tipo_categorias))
	frecuencias[5] = make([]int, len(cl.Sub_tipo_categorias))
	//usamos arreglos paralelos, ya tenemos indexados todos los valores unicos de cada elemento
	for i := 0; i < k; i++ {
		frecuencias[0][searchindex(cl.Regiones, (*distancias[k].destino).REGIÓN)]++
		frecuencias[1][searchindex(cl.Provincias, (*distancias[k].destino).PROVINCIA)]++
		frecuencias[2][searchindex(cl.Distritos, (*distancias[k].destino).DISTRITO)]++
		frecuencias[3][searchindex(cl.Categorias, (*distancias[k].destino).CATEGORIA)]++
		frecuencias[4][searchindex(cl.Tipo_categorias, (*distancias[k].destino).Tipo_de_Categoria)]++
		frecuencias[5][searchindex(cl.Sub_tipo_categorias, (*distancias[k].destino).Sub_tipo_Categoria)]++
	}
	recursoKnn <- cl.RecursoKnn{
		Recurso:                    &origen,
		REGIÓN_predict:             cl.Regiones[getindicemayor(frecuencias[0])],
		PROVINCIA_predict:          cl.Provincias[getindicemayor(frecuencias[1])],
		DISTRITO_predict:           cl.Distritos[getindicemayor(frecuencias[2])],
		CATEGORIA_predict:          cl.Categorias[getindicemayor(frecuencias[3])],
		Tipo_de_Categoria_predict:  cl.Tipo_categorias[getindicemayor(frecuencias[4])],
		Sub_tipo_Categoria_predict: cl.Sub_tipo_categorias[getindicemayor(frecuencias[5])],
	}
}

//hallamos las distancias de 1 punto a todos (todos con todos)
//para cada punto de sus distancias solo agarra los k mas cercanos
//de los cercanos hallamos sus respectivas categorias
func Knn(k int, recursos []cl.Recurso) []cl.RecursoKnn {
	if k > len(cl.Recursos) {
		k = len(cl.Recursos)
	}
	channel := make(chan cl.RecursoKnn)
	//creamos canales para llenar la matriz
	cant := 0
	for i := 0; i < len(recursos); i++ {
		if recursos[i].LATITUD != nil && recursos[i].LONGITUD != nil {
			go predictRecurso(recursos[i], k, channel)
			cant++
		}
	}
	//leemos del canal:
	var aux []cl.RecursoKnn
	for recurso := range channel {
		aux = append(aux, recurso)
		cant--
		if cant == 0 {
			close(channel)
		}
	}
	return aux
}
func appendIfMissing(strs []string, str string) []string {
	for _, aux := range strs {
		if aux == str {
			return strs
		}
	}
	return append(strs, str)
}
func GetUniqueValues() {
	for _, recurso := range cl.Recursos {
		cl.Regiones = appendIfMissing(cl.Regiones, recurso.REGIÓN)
		cl.Provincias = appendIfMissing(cl.Provincias, recurso.PROVINCIA)
		cl.Distritos = appendIfMissing(cl.Distritos, recurso.DISTRITO)
		cl.Categorias = appendIfMissing(cl.Categorias, recurso.CATEGORIA)
		cl.Tipo_categorias = appendIfMissing(cl.Tipo_categorias, recurso.Tipo_de_Categoria)
		cl.Sub_tipo_categorias = appendIfMissing(cl.Sub_tipo_categorias, recurso.Sub_tipo_Categoria)
	}
}
