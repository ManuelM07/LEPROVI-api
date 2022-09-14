package main

import (
	"fmt"
)

/*
Esta función se encarga de dar inicio al parsing, dependiendo el tipo de nodo, hace el llamado a la función
correspondiente y concatenando su resultado en la variable de tipo string code, finalmente está función
retorna la variable code, que contiene el codigo formado apartir de los nodos.
*/
func startParsingJs() string {
	var countMath = map[string]int{"add": 0, "less": 0, "mult": 0, "divide": 0, "module": 0}
	var code string

	for k := 0; k < len(nodes); k++ {
		if nodes[k].name == "NodeMath" {
			methodNode := fmt.Sprintf("%v", nodes[k].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
			if countMath[methodNode] == 0 { // esta condicción se hace con el fin de controlar las funciones Math, si una función ya fue creada, no es necesario volver a crearla.
				code += mathOperationJs("body", k)
				countMath[methodNode] = 1
			}
		} else if nodes[k].name == "NodeAssign" {
			code += assignJs(k, nodes[k].inputs.(map[string]interface{}))
		} else if nodes[k].name == "NodePrint" {
			code += printJs(nodes[k].inputs.(map[string]interface{}))
		} else if nodes[k].name == "NodeIf" {
			code += nodeIfJs(k)
		} else if nodes[k].name == "NodeElse" {
			code += nodeElseJs(k)
		} else if nodes[k].name == "NodeFor" {
			code += nodeForJs(k)
		}
	}
	return code

}

func mathOperationJs(option string, idOutput int) string {
	operation := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idOutput].inputs.(map[string]interface{})
	if option == "body" {
		return fmt.Sprintf("function %s(a, b) {\n\treturn a%sb\n}\n", operation, typeOperation[operation])
	} else {
		node1 := findInput(inputs["input_1"])
		node2 := findInput(inputs["input_2"])
		nodePos1 := findNode(node1, nodes)
		nodePos2 := findNode(node2, nodes)
		return fmt.Sprintf("%s(%s, %s)", operation, typeNodeJs(nodes[nodePos1], nodePos1), typeNodeJs(nodes[nodePos2], nodePos2))
	}
}

func assignJs(pos int, inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1, nodes)
	varName := nodes[pos].data.(map[string]interface{})["url"]
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperationJs("", idNode)
		return fmt.Sprintf("%s = %s\n", varName, answer)
	} else if nodes[idNode].name == "NodeIf" || nodes[idNode].name == "NodeElse" || nodes[idNode].name == "NodeFor" {
		answer := valueAssigned(idNode) // Para el caso del else, se puede reutilizar la funcion de nodeIf
		return fmt.Sprintf("\t%s = %s\n}\n", varName, answer)
	} else if nodes[idNode].name == "NodeNumber" || nodes[idNode].name == "NodeString" || nodes[idNode].name == "NodeAssign" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("%s = %s\n", varName, answer)
	} else if nodes[idNode].name == "NodeStringOp" {
		answer := stringOperationsJs(idNode)
		return fmt.Sprintf("%s = %s\n", varName, answer)
	}
	return ""
}

func printJs(inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1, nodes)
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("n", idNode)
		return fmt.Sprintf("\nconsole.log(%s)\n", answer)
	} else if nodes[idNode].name == "NodeAssign" || nodes[idNode].name == "NodeNumber" || nodes[idNode].name == "NodeString" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("console.log(%s)\n", answer)
	} else if nodes[idNode].name == "NodeIf" || nodes[idNode].name == "NodeElse" || nodes[idNode].name == "NodeFor" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("\tconsole.log(%s)\n}\n", answer)
	} else if nodes[idNode].name == "NodeStringOp" {
		return fmt.Sprintf("console.log(%s)\n", stringOperationsJs(idNode))
	}
	return ""
}

func nodeIfJs(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	node3 := findInput(inputs["input_3"])
	nodePos1 := findNode(node1, nodes)
	nodePos2 := findNode(node2, nodes)
	nodePos3 := findNode(node3, nodes)
	return fmt.Sprintf("if (%s %s %s) {\n", typeNodeJs(nodes[nodePos1], nodePos1), typeNodeJs(nodes[nodePos2], nodePos2), typeNodeJs(nodes[nodePos3], nodePos3))
}

func nodeElseJs(idNode int) string {
	return "else {\n"
}

func nodeForJs(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	nodePos1 := findNode(node1, nodes)
	nodePos2 := findNode(node2, nodes)
	return fmt.Sprintf("for (let index = %s; index < %s; index++) {\n", typeNodeJs(nodes[nodePos1], nodePos1), typeNodeJs(nodes[nodePos2], nodePos2))
}

func comparisonJs(idOutput int) string {
	comparison := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	return typeComparison[comparison]
}

func stringOperationsJs(idNode int) string {
	operation := fmt.Sprintf("%v", nodes[idNode].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	nodePos1 := findNode(node1, nodes)
	if operation == "len" {
		return fmt.Sprintf("%v.length", typeNodeJs(nodes[nodePos1], nodePos1))
	} else if operation == "first" {
		return fmt.Sprintf("%v.slice(0, 1)", typeNodeJs(nodes[nodePos1], nodePos1))
	} else if operation == "rest" {
		return fmt.Sprintf("%v.slice(1)", typeNodeJs(nodes[nodePos1], nodePos1))
	}
	return ""
}

//--------------------------- Funciones auxiliares ---------------------------\\

/*
Esta función se encarga de recibir un nodo y retornar su respuesta, dependiendo el tipo de nodo
si es un nodo Number, retorna el numero asociado al nodo, en caso de que sea un nodo Assign,
retorna la variable asociada al mismo, en caso de comparison o mathOperation se llama a la función
correspondiente.
*/
func typeNodeJs(node Node, posNode int) string {
	nameNode := fmt.Sprintf("%v", node.name)
	if nameNode == "NodeNumber" || nameNode == "NodeAssign" || nameNode == "NodeString" {
		return valueAssigned(posNode)
	} else if nameNode == "NodeComOp" {
		return comparison(posNode)
	} else if nameNode == "NodeMath" {
		return mathOperationJs("", posNode)
	} else if nameNode == "NodeStringOp" {
		return stringOperationsJs(posNode)
	}
	return ""
}
