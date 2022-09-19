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
		} /*else if nodes[k].name == "NodeNumber" {
			code += nodeNumber(k)
		}*/
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
			return fmt.Sprintf("%s%s = %s", prompter, varName, answer)
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

//--------------------------- Funciones auxiliares ---------------------------\\

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

/*
Esta función se encarga de ordenar los nodos. Siendo C un conjunto de nodos,
C1, será el nodo que inicia las relación entre los nodos del conjunto y CN
el nodo que las finaliza.
Fanalmente retorna el conjunto de nodos ordenados.
*/
/*func sortNodes(nodeAux []Node) []Node {
	var sort []Node
	var posI int
	var isIf bool
	var father bool     // se usa para identificar si es un C1
	var brothers []Node // es un slice de Node, donde se guardaran los hermanos para cada iteración
	n := -1             // n será una variable acumuladora, por cada nodo que se agregue a sort(Node), esta aumentará en 1
	var auxOutputs bool

	for {
		fmt.Println("Nodes:", sort)
		if len(nodeAux) == 1 { // Si queda un solo nodo, este se agrega y se rompe el ciclo
			sort = append(sort, nodeAux[0])
			break
		} else if isIf && nodeAux[posI].name == "NodeElse" || !isIf && len(nodeAux[posI].inputs.(map[string]interface{})["input_1"].(map[string]interface{})["connections"].([]interface{})) == 0 || auxOutputs { // se busca el padre del arbol
			//if nodeAux[posI].name == "NodeAssign"
			auxOutputs = false
			father = true
			isIf = false
			sort = append(sort, nodeAux[posI])
			nodeAux = RemoveIndex(nodeAux, posI)
			n++

			thisNode := sort[n].outputs.(map[string]interface{})["output_1"]
			for {
				brothers = nil
				if thisNode != nil { // Para tener un hermano, antes debe tener una salida
					if len(thisNode.(map[string]interface{})["connections"].([]interface{})) == 0 { // si no tiene conexiones es porque es un nodo print, en este caso se cierra el ciclo
						break
					}
					thisOutput := findInput(thisNode)[0]
					for i := 0; i < len(nodeAux); i++ { // se busca si ese padre tiene un hermano
						if nodeAux[i].outputs.(map[string]interface{})["output_1"] == nil || len(nodeAux[i].outputs.(map[string]interface{})["output_1"].(map[string]interface{})["connections"].([]interface{})) != 0 {
							if nodeAux[i].outputs.(map[string]interface{})["output_1"] != nil && len(nodeAux[i].outputs.(map[string]interface{})["output_1"].(map[string]interface{})["connections"].([]interface{})) > 1 {
								auxOutputs = true
							}

							nodesInpunts := findInput(nodeAux[i].outputs.(map[string]interface{})["output_1"])
							for p := 0; p < len(nodesInpunts); p++ {
								if nodeAux[i].name != "NodePrint" && nodesInpunts[p] == thisOutput || len(nodesInpunts) > 1 { // se verifica si la salida de X nodo es igual a la del nodo padre, si se cumple entonces son hermanos
									if nodeAux[i].inputs == nil || len(nodesInpunts) > 1 { // si no tiene input, esto quiere decir que es un hermano que no tiene padre, por lo anterior es un C1
										brothers = append(brothers, nodeAux[i])
									} else { // si tiene padre, esto implica que no es un C1, por lo anterior se rompe el ciclo y se sigue buscando el C1
										brothers = nil
										father = false
										break
									}
									for j := 0; j < len(brothers); j++ {
										fmt.Println("Bros: ", brothers)
										sort = append(sort, brothers[j])
										nodeAux = RemoveIndex(nodeAux, findNode(brothers[j].id, nodeAux))
										n++
									}
								}
							}
						}
					}
					if !father { // como no es C1 se rompe el otro ciclo
						break
					}

					childNodePos := findNode(thisOutput, nodeAux)
					if childNodePos == -1 { // si se obtiene como resultado -1, esto nos indica que ya no hay mas relaciones en ese conjuto de nodos, por lo anterior se rompe el ciclo
						break
					} else if nodeAux[childNodePos].name == "NodeIf" { // se verifica si es un nodoIf
						isIf = true
					}
					thisNode = nodeAux[childNodePos].outputs.(map[string]interface{})["output_1"]
					sort = append(sort, nodeAux[childNodePos])
					nodeAux = RemoveIndex(nodeAux, childNodePos)
					n++
				} else { // si el nodo prensente en thisNode es nulo, entonces se rompe el ciclo
					break
				}
			}
			posI = -1
		}
		if len(nodeAux) == 0 {
			break
		}
		posI++
	}
	return sort
}*/

/*
Esta función se encarga de eliminar un elemento de una lista.
*/
/*func RemoveIndex(s []Node, index int) []Node {
	if len(s) == 1 {
		return nil
	}
	return append(s[:index], s[index+1:]...)
}*/
