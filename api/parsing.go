package main

import (
	"encoding/json"
	"fmt"
)

var nodes []Node

type Node struct {
	id      string
	name    interface{}
	inputs  interface{}
	outputs interface{}
	data    interface{}
}

func mapJson(data string) {
	nodes = nil

	var m map[string]interface{}
	json.Unmarshal([]byte(data), &m)

	for key, element := range m {

		node := Node{
			id:      key,
			name:    element.(map[string]interface{})["name"],
			inputs:  element.(map[string]interface{})["inputs"],
			outputs: element.(map[string]interface{})["outputs"],
			data:    element.(map[string]interface{})["data"],
		}
		nodes = append(nodes, node)
	}
	startParsing()
}

func startParsing() {
	for k := 0; k < len(nodes); k++ {
		if nodes[k].name == "NodeMath" {
			//fmt.Println(nodes[k].inputs)
			fmt.Println(mathOperation("body", k)) // para acceder a la clave de un map de varios niveles
		} else if nodes[k].name == "NodeAssing" {
			fmt.Println(assing(k, nodes[k].inputs.(map[string]interface{})))
		} else if nodes[k].name == "NodePrint" {
			fmt.Println(print(nodes[k].inputs.(map[string]interface{})))
			//fmt.Println(findInput(nodes[k].inputs.(map[string]interface{})["input_1"]))
		}
	}

}

func mathOperation(option string, idOutput int) string {
	operation := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idOutput].inputs.(map[string]interface{})
	if option == "body" {
		if operation == "add" {
			return "def add(a, b):\n\treturn a+b\n"
		} else if operation == "less" {
			return "def less(a, b):\n\treturn a-b\n"
		}
	} else {
		node1 := findInput(inputs["input_1"])
		node2 := findInput(inputs["input_2"])
		return fmt.Sprintf("%s(%s, %s)", operation, nodes[findNode(node1)].data.(map[string]interface{})["url"], nodes[findNode(node2)].data.(map[string]interface{})["url"])
	}
	return ""
}

func assing(pos int, inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1)
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("n", idNode)
		varName := nodes[pos].data.(map[string]interface{})["url"]
		return fmt.Sprintf("%s = %s", varName, answer)
	}
	return ""
}

func print(inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1)
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("n", idNode)
		return fmt.Sprintf("print(%s)", answer)
	} else if nodes[idNode].name == "NodeAssing" {
		varName := nodes[idNode].data.(map[string]interface{})["url"]
		return fmt.Sprintf("print(%s)", varName)
	}
	return ""
}

//--------------------------- Funciones auxiliares ---------------------------\\

/**
* Esta función se encarga de buscar un nodo en un array(slice) de nodos
 */
func findNode(id string) int {
	for k := 0; k < len(nodes); k++ {
		if nodes[k].id == id {
			return k
		}
	}
	return -1
}

/**
* Esta función se encarga de buscar un nodo input en una interface de inputs
 */
func findInput(input interface{}) string {
	return fmt.Sprintf("%v", input.(map[string]interface{})["connections"].([]interface{})[0].(map[string]interface{})["node"]) // fmt.Sprintf("%v", node1) permite convertir una interfaz en string
}
