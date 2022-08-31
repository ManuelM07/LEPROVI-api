package main

import (
	"encoding/json"
	"fmt"
)

var nodes []Node
var typeOperation = map[string]string{"add": "+", "less": "-"} // se guardan todos los tipos de operaciones matematicas
var typeComparison = map[string]string{"equals": "==", "greater": ">"}

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
		} else if nodes[k].name == "NodeIf" {
			fmt.Println(nodeIf("body", k))
		}
	}

}

func mathOperation(option string, idOutput int) string {
	operation := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idOutput].inputs.(map[string]interface{})
	if option == "body" {
		return fmt.Sprintf("def %s(a, b):\n\treturn a%sb\n", operation, typeOperation[operation])
	} else {
		node1 := findInput(inputs["input_1"])
		node2 := findInput(inputs["input_2"])
		nodePos1 := findNode(node1)
		nodePos2 := findNode(node2)
		return fmt.Sprintf("%s(%s, %s)", operation, typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2))
	}
}

func assing(pos int, inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1)
	varName := nodes[pos].data.(map[string]interface{})["url"]
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("", idNode)
		return fmt.Sprintf("%s = %s", varName, answer)
	} else if nodes[idNode].name == "NodeIf" {
		answer := nodeIf("", idNode)
		return fmt.Sprintf("\t%s = %s", varName, answer)
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

func nodeIf(option string, idNode int) string {
	if option == "body" {
		inputs := nodes[idNode].inputs.(map[string]interface{})
		node1 := findInput(inputs["input_1"])
		node2 := findInput(inputs["input_2"])
		node3 := findInput(inputs["input_3"])
		nodePos1 := findNode(node1)
		nodePos2 := findNode(node2)
		nodePos3 := findNode(node3)
		return fmt.Sprintf("if %s %s %s:\n", typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2), typeNode(nodes[nodePos3], nodePos3))
	} else {
		return fmt.Sprintf("%v", nodes[idNode].data.(map[string]interface{})["url"])
	}
}

func comparison(idOutput int) string {
	comparison := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	return typeComparison[comparison]
}

//--------------------------- Funciones auxiliares ---------------------------\\

/**
* Esta función se encarga de buscar la posición de un nodo en un array(slice) de nodos
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

/**
* Esta función se encarga de recibir un nodo y retornar su respuesta, dependiendo el tipo de nodo
 */
func typeNode(node Node, posNode int) string {
	nameNode := fmt.Sprintf("%v", node.name)
	if nameNode == "NodeNumber" {
		return fmt.Sprintf("%v", node.data.(map[string]interface{})["url"])
	} else if nameNode == "NodeComOp" {
		return comparison(posNode)
	}
	return ""
}
