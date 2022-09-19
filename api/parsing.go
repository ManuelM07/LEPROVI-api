package main

import (
	"fmt"
	"sort"
)

var nodes []Node
var typeOperation = map[string]string{"add": "+", "less": "-", "mult": "*", "divide": "/", "module": "%"}                                   // se guardan todos los tipos de operaciones matematicas
var typeComparison = map[string]string{"equals": "==", "greater": ">", "less": "<", "greaterOrE": ">=", "lessOrE": "<=", "different": "!="} // se guardan todos los operadores de comparación
var prompter string

type Node struct {
	id      string
	name    interface{}
	inputs  interface{}
	outputs interface{}
	data    interface{}
	posX    float64
}

/*
Esta función se encarga de mapear el json(este antes fue convertido a string) que recibe desde el front,
sacando los elementos mas importantes y guardarlos en una estructura Node.
*/
func mapJson(data map[string]interface{}, languaje string) string {
	nodes = nil

	for key, element := range data {

		node := Node{
			id:      key,
			name:    element.(map[string]interface{})["name"],
			inputs:  element.(map[string]interface{})["inputs"],
			outputs: element.(map[string]interface{})["outputs"],
			data:    element.(map[string]interface{})["data"],
			posX:    element.(map[string]interface{})["pos_x"].(float64),
		}
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].posX < nodes[j].posX }) // Permite ordenar los nodos por su posición en X, esto es util cuando se tiene mas de un bloque o conjunto de nodos
	//nodes = sortNodes(nodes)
	if languaje == "nodejs" {
		return startParsingJs()
	}
	return startParsing()
}

/*
Esta función se encarga de dar inicio al parsing, dependiendo el tipo de nodo, hace el llamado a la función
correspondiente y concatenando su resultado en la variable de tipo string code, finalmente está función
retorna la variable code, que contiene el codigo formado apartir de los nodos.
*/
func startParsing() string {
	var countMath = map[string]int{"add": 0, "less": 0, "mult": 0, "divide": 0, "module": 0}
	var code string
	prompter = ""

	for k := 0; k < len(nodes); k++ {
		if nodes[k].inputs.(map[string]interface{})["input_1"] == nil || len(nodes[k].inputs.(map[string]interface{})["input_1"].(map[string]interface{})["connections"].([]interface{})) == 0 {
			if nodes[k].name != "NodeMoveLeft" {
				prompter = ""
			}
		}
		if len(prompter) > 0 {
			if nodes[k].name == "NodeMoveLeft" {
				prompter = string(prompter[0 : len(prompter)-1]) // elimina una tabulación
			}
		}

		if nodes[k].name == "NodeMath" {
			methodNode := fmt.Sprintf("%v", nodes[k].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
			if countMath[methodNode] == 0 { // esta condicción se hace con el fin de controlar las funciones Math, si una función ya fue creada, no es necesario volver a crearla.
				code += mathOperation("body", k)
				countMath[methodNode] = 1
			}
		} else if nodes[k].name == "NodeAssign" {
			code += assign(k, nodes[k].inputs.(map[string]interface{}))
		} else if nodes[k].name == "NodePrint" {
			code += print(nodes[k].inputs.(map[string]interface{}))
		} else if nodes[k].name == "NodeIf" {
			code += nodeIf(k)
			prompter += "\t"
		} else if nodes[k].name == "NodeElse" {
			code += nodeElse(k)
			prompter += "\t"
		} else if nodes[k].name == "NodeFor" {
			code += nodeFor(k)
			prompter += "\t"
		}
	}
	return code
}

func mathOperation(option string, idOutput int) string {
	operation := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idOutput].inputs.(map[string]interface{})
	if option == "body" {
		return fmt.Sprintf("%[1]sdef %[2]s(a, b):\n\t%[1]sreturn a%[3]sb\n", prompter, operation, typeOperation[operation])
	} else {
		node1 := findInput(inputs["input_1"])
		node2 := findInput(inputs["input_2"])
		nodePos1 := findNode(node1[0], nodes)
		nodePos2 := findNode(node2[0], nodes)
		return fmt.Sprintf("%s(%s, %s)", operation, typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2))
	}
}

func assign(pos int, inputs map[string]interface{}) string {
	if len(inputs["input_1"].(map[string]interface{})["connections"].([]interface{})) != 0 {
		node1 := findInput(inputs["input_1"])
		idNode := findNode(node1[0], nodes)
		varName := nodes[pos].data.(map[string]interface{})["url"]
		if nodes[idNode].name == "NodeMath" {
			answer := mathOperation("", idNode)
			return fmt.Sprintf("%s%s = %s\n", prompter, varName, answer)
		} else if nodes[idNode].name == "NodeNumber" || nodes[idNode].name == "NodeString" || nodes[idNode].name == "NodeAssign" {
			answer := valueAssigned(idNode)
			return fmt.Sprintf("%s%s = %s\n", prompter, varName, answer)
		} else if nodes[idNode].name == "NodeStringOp" {
			answer := stringOperations(idNode)
			return fmt.Sprintf("%s%s = %s", prompter, varName, answer)
		}
		/*else if nodes[idNode].name == "NodeIf" || nodes[idNode].name == "NodeElse" {
			answer := valueAssigned(idNode) // Para el caso del else y for, se puede reutilizar la funcion de nodeIf
			return fmt.Sprintf("%s\t%s = %s\n", prompter, varName, answer)
		}*/
	}
	return ""
}

func print(inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1[0], nodes)
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("n", idNode)
		return fmt.Sprintf("%sprint(%s)\n", prompter, answer)
	} else if nodes[idNode].name == "NodeAssign" || nodes[idNode].name == "NodeNumber" || nodes[idNode].name == "NodeString" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("%sprint(%s)\n", prompter, answer)
	} else if nodes[idNode].name == "NodeIf" || nodes[idNode].name == "NodeElse" || nodes[idNode].name == "NodeFor" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("%sprint(%s)\n", prompter, answer)
	} else if nodes[idNode].name == "NodeStringOp" {
		return fmt.Sprintf("%sprint(%s)\n", prompter, stringOperations(idNode))
	}
	return ""
}

func nodeIf(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	node3 := findInput(inputs["input_3"])
	nodePos1 := findNode(node1[0], nodes)
	nodePos2 := findNode(node2[0], nodes)
	nodePos3 := findNode(node3[0], nodes)
	return fmt.Sprintf("%sif %s %s %s:\n", prompter, typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2), typeNode(nodes[nodePos3], nodePos3))
}

func nodeElse(idNode int) string {
	return fmt.Sprintf("%selse:\n", prompter)
}

func nodeFor(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	nodePos1 := findNode(node1[0], nodes)
	nodePos2 := findNode(node2[0], nodes)
	return fmt.Sprintf("%sfor i in range(%s, %s):\n", prompter, typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2))
}

func comparison(idOutput int) string {
	comparison := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	return typeComparison[comparison]
}

func stringOperations(idNode int) string {
	operation := fmt.Sprintf("%v", nodes[idNode].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	nodePos1 := findNode(node1[0], nodes)
	if operation == "len" {
		return fmt.Sprintf("len(%v)", typeNode(nodes[nodePos1], nodePos1))
	} else if operation == "first" {
		return fmt.Sprintf("%v[0]", typeNode(nodes[nodePos1], nodePos1))
	} else if operation == "rest" {
		return fmt.Sprintf("%v[1:]", typeNode(nodes[nodePos1], nodePos1))
	}
	return ""
}

//--------------------------- Funciones auxiliares ---------------------------//

/*
Esta función se encarga de buscar la posición de un nodo en un array(slice) de nodos.
*/
func findNode(id string, nodesX []Node) int {
	for k := 0; k < len(nodesX); k++ {
		if nodesX[k].id == id {
			return k
		}
	}
	return -1
}

/*
Esta función se encarga de buscar un nodo input en una interface de inputs, a su vez también sirve para buscar un nodo outputs.
*/
func findInput(input interface{}) []string {
	var sliceInputs []string
	inputs := input.(map[string]interface{})["connections"].([]interface{})
	for k := 0; k < len(inputs); k++ {
		sliceInputs = append(sliceInputs, fmt.Sprintf("%v", inputs[k].(map[string]interface{})["node"]))
	}
	return sliceInputs //fmt.Sprintf("%v", inputs.(map[string]interface{})["node"]) // fmt.Sprintf("%v", node1) permite convertir una interfaz en string
}

/*
Esta función se encarga de recibir un nodo y retornar su respuesta, dependiendo el tipo de nodo
si es un nodo Number, retorna el numero asociado al nodo, en caso de que sea un nodo Assign,
retorna la variable asociada al mismo, en caso de comparison o mathOperation se llama a la función
correspondiente.
*/
func typeNode(node Node, posNode int) string {
	nameNode := fmt.Sprintf("%v", node.name)
	if nameNode == "NodeNumber" || nameNode == "NodeAssign" || nameNode == "NodeString" {
		return valueAssigned(posNode)
	} else if nameNode == "NodeComOp" {
		return comparison(posNode)
	} else if nameNode == "NodeMath" {
		return mathOperation("", posNode)
	} else if nameNode == "NodeStringOp" {
		return stringOperations(posNode)
	}
	return ""
}

/*
Esta función se encarga de retornar el valor a ser asignado en la funcion assign.
*/
func valueAssigned(posNode int) string {
	if nodes[posNode].name == "NodeString" {
		return fmt.Sprintf("'%v'", nodes[posNode].data.(map[string]interface{})["url"])
	}
	return fmt.Sprintf("%v", nodes[posNode].data.(map[string]interface{})["url"])
}
