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
	var closeKey int // contador, closeKey indica la cantidad de llaves({) que se deben cerrar
	prompter = ""

	for k := 0; k < len(nodes); k++ {
		if nodes[k].inputs.(map[string]interface{})["input_1"] == nil || len(nodes[k].inputs.(map[string]interface{})["input_1"].(map[string]interface{})["connections"].([]interface{})) == 0 {
			if nodes[k].name != "NodeMoveLeft" {
				if closeKey > 0 {
					prompter = string(prompter[0 : len(prompter)-1]) // elimina una tabulación
					code += fmt.Sprintf("%s}\n", prompter)
					closeKey -= 1
				}
				prompter = ""
			}
		}
		if len(prompter) > 0 {
			if nodes[k].name == "NodeMoveLeft" {
				prompter = string(prompter[0 : len(prompter)-1])
				if closeKey > 0 {
					code += fmt.Sprintf("%s}\n", prompter)
					closeKey -= 1
				}
			}
		}

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
			prompter += "\t"
			closeKey += 1
		} else if nodes[k].name == "NodeElse" {
			code += nodeElseJs(k)
			prompter += "\t"
			closeKey += 1
		} else if nodes[k].name == "NodeFor" {
			code += nodeForJs(k)
			prompter += "\t"
			closeKey += 1
		}
	}
	if closeKey > 0 { // en caso tal de que no se hayan cerrado todas las llaves, se procede a cerrarlas por medio del for
		for k := 0; k < closeKey; k++ {
			if len(prompter) != 0 {
				prompter = string(prompter[0 : len(prompter)-1])
			}
			code += fmt.Sprintf("%s}\n", prompter)
		}
	}
	return code

}

func mathOperationJs(option string, idOutput int) string {
	operation := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idOutput].inputs.(map[string]interface{})
	if option == "body" {
		return fmt.Sprintf("%[1]sfunction %[2]s(a, b) {\n\t%[1]sreturn a%[3]sb\n%[1]s}\n", prompter, operation, typeOperation[operation])
	} else {
		node1 := findInput(inputs["input_1"])
		node2 := findInput(inputs["input_2"])
		nodePos1 := findNode(node1[0], nodes)
		nodePos2 := findNode(node2[0], nodes)
		return fmt.Sprintf("%s(%s, %s)", operation, typeNodeJs(nodes[nodePos1], nodePos1), typeNodeJs(nodes[nodePos2], nodePos2))
	}
}

func assignJs(pos int, inputs map[string]interface{}) string {
	if len(inputs["input_1"].(map[string]interface{})["connections"].([]interface{})) != 0 {
		node1 := findInput(inputs["input_1"])
		idNode := findNode(node1[0], nodes)
		varName := nodes[pos].data.(map[string]interface{})["url"]
		if nodes[idNode].name == "NodeMath" {
			answer := mathOperationJs("", idNode)
			return fmt.Sprintf("%s%s = %s\n", prompter, varName, answer)
		} else if nodes[idNode].name == "NodeNumber" || nodes[idNode].name == "NodeString" || nodes[idNode].name == "NodeAssign" {
			answer := valueAssigned(idNode)
			return fmt.Sprintf("%s%s = %s\n", prompter, varName, answer)
		} else if nodes[idNode].name == "NodeStringOp" {
			answer := stringOperationsJs(idNode)
			return fmt.Sprintf("%s%s = %s\n", prompter, varName, answer)
		}
	}
	return ""
}

func printJs(inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1[0], nodes)
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("n", idNode)
		return fmt.Sprintf("%s\nconsole.log(%s)\n", prompter, answer)
	} else if nodes[idNode].name == "NodeAssign" || nodes[idNode].name == "NodeNumber" || nodes[idNode].name == "NodeString" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("%sconsole.log(%s)\n", prompter, answer)
	} else if nodes[idNode].name == "NodeIf" || nodes[idNode].name == "NodeElse" || nodes[idNode].name == "NodeFor" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("%s\tconsole.log(%s)\n}\n", prompter, answer)
	} else if nodes[idNode].name == "NodeStringOp" {
		return fmt.Sprintf("%sconsole.log(%s)\n", prompter, stringOperationsJs(idNode))
	}
	return ""
}

func nodeIfJs(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	node3 := findInput(inputs["input_3"])
	nodePos1 := findNode(node1[0], nodes)
	nodePos2 := findNode(node2[0], nodes)
	nodePos3 := findNode(node3[0], nodes)
	return fmt.Sprintf("%sif (%s %s %s) {\n", prompter, typeNodeJs(nodes[nodePos1], nodePos1), typeNodeJs(nodes[nodePos2], nodePos2), typeNodeJs(nodes[nodePos3], nodePos3))
}

func nodeElseJs(idNode int) string {
	return fmt.Sprintf("%selse {\n", prompter)
}

func nodeForJs(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	nodePos1 := findNode(node1[0], nodes)
	nodePos2 := findNode(node2[0], nodes)
	return fmt.Sprintf("%sfor (let index = %s; index < %s; index++) {\n", prompter, typeNodeJs(nodes[nodePos1], nodePos1), typeNodeJs(nodes[nodePos2], nodePos2))
}

func stringOperationsJs(idNode int) string {
	operation := fmt.Sprintf("%v", nodes[idNode].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	nodePos1 := findNode(node1[0], nodes)
	if operation == "len" {
		return fmt.Sprintf("%v.length", typeNodeJs(nodes[nodePos1], nodePos1))
	} else if operation == "first" {
		return fmt.Sprintf("%v.slice(0, 1)", typeNodeJs(nodes[nodePos1], nodePos1))
	} else if operation == "rest" {
		return fmt.Sprintf("%v.slice(1)", typeNodeJs(nodes[nodePos1], nodePos1))
	}
	return ""
}

//--------------------------- Funciones auxiliares ---------------------------//

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
