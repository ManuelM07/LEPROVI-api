package main

import (
	"encoding/json"
	"fmt"
)

var nodes []Node
var typeOperation = map[string]string{"add": "+", "less": "-"} // se guardan todos los tipos de operaciones matematicas
var typeComparison = map[string]string{"equals": "==", "greater": ">"}

//var tree map[string]interface{}

type Node struct {
	id      string
	name    interface{}
	inputs  interface{}
	outputs interface{}
	data    interface{}
	//father 	interface{}
}

/*
Esta función se encarga de mapear el json(este antes fue convertido a string) que recibe desde el front,
sacando los elementos mas importantes y guardarlos en una estructura Node.
*/
func mapJson(data string) string {
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
	return startParsing()
}

/*
Esta función se encarga de dar inicio al parsing, dependiendo el tipo de nodo, hace el llamado a la función
correspondiente y concatenando su resultado en la variable de tipo string code, finalmente está función
retorna la variable code, que contiene el codigo formado apartir de los nodos.
*/
func startParsing() string {
	nodes = sortNodes(nodes)
	var code string
	//fmt.Println((nodes))
	for k := 0; k < len(nodes); k++ {
		if nodes[k].name == "NodeMath" {
			code += mathOperation("body", k) // para acceder a la clave de un map de varios niveles
		} else if nodes[k].name == "NodeAssign" {
			code += assign(k, nodes[k].inputs.(map[string]interface{}))
		} else if nodes[k].name == "NodePrint" {
			code += print(nodes[k].inputs.(map[string]interface{}))
			//fmt.Println(findInput(nodes[k].inputs.(map[string]interface{})["input_1"]))
		} else if nodes[k].name == "NodeIf" {
			code += nodeIf(k)
		} else if nodes[k].name == "NodeElse" {
			code += nodeElse(k)
		} else if nodes[k].name == "NodeFor" {
			code += nodeFor(k)
		}
	}
	return code
	//fmt.Println(sortNodes(nodes))

}

func mathOperation(option string, idOutput int) string {
	operation := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	inputs := nodes[idOutput].inputs.(map[string]interface{})
	if option == "body" {
		return fmt.Sprintf("def %s(a, b):\n\treturn a%sb\n", operation, typeOperation[operation])
	} else {
		node1 := findInput(inputs["input_1"])
		node2 := findInput(inputs["input_2"])
		nodePos1 := findNode(node1, nodes)
		nodePos2 := findNode(node2, nodes)
		return fmt.Sprintf("%s(%s, %s)", operation, typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2))
	}
}

func assign(pos int, inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1, nodes)
	varName := nodes[pos].data.(map[string]interface{})["url"]
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("", idNode)
		return fmt.Sprintf("%s = %s\n", varName, answer)
	} else if nodes[idNode].name == "NodeIf" || nodes[idNode].name == "NodeElse" || nodes[idNode].name == "NodeFor" {
		answer := valueAssigned(idNode) // Para el caso del else, se puede reutilizar la funcion de nodeIf
		return fmt.Sprintf("\t%s = %s\n", varName, answer)
	} else if nodes[idNode].name == "NodeNumber" {
		answer := valueAssigned(idNode)
		return fmt.Sprintf("%s = %s\n", varName, answer)
	}
	return ""
}

func print(inputs map[string]interface{}) string {
	node1 := findInput(inputs["input_1"])
	idNode := findNode(node1, nodes)
	if nodes[idNode].name == "NodeMath" {
		answer := mathOperation("n", idNode)
		return fmt.Sprintf("print(%s)\n", answer)
	} else if nodes[idNode].name == "NodeAssign" || nodes[idNode].name == "NodeNumber" {
		varName := valueAssigned(idNode)
		return fmt.Sprintf("print(%s)\n", varName)
	} else if nodes[idNode].name == "NodeFor" {
		varName := valueAssigned(idNode)
		return fmt.Sprintf("\tprint(%s)\n", varName)
	}
	return ""
}

func nodeIf(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	node3 := findInput(inputs["input_3"])
	nodePos1 := findNode(node1, nodes)
	nodePos2 := findNode(node2, nodes)
	nodePos3 := findNode(node3, nodes)
	return fmt.Sprintf("if %s %s %s:\n", typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2), typeNode(nodes[nodePos3], nodePos3))
}

func nodeElse(idNode int) string {
	return "else:\n"
}

func nodeFor(idNode int) string {
	inputs := nodes[idNode].inputs.(map[string]interface{})
	node1 := findInput(inputs["input_1"])
	node2 := findInput(inputs["input_2"])
	nodePos1 := findNode(node1, nodes)
	nodePos2 := findNode(node2, nodes)
	return fmt.Sprintf("for i in range(%s, %s):\n", typeNode(nodes[nodePos1], nodePos1), typeNode(nodes[nodePos2], nodePos2))
}

func comparison(idOutput int) string {
	comparison := fmt.Sprintf("%v", nodes[idOutput].data.(map[string]interface{})["data"].(map[string]interface{})["method"])
	return typeComparison[comparison]
}

//--------------------------- Funciones auxiliares ---------------------------\\

/*
Esta función se encarga de buscar la posición de un nodo en un array(slice) de nodos.
*/
func findNode(id string, nodesX []Node) int {
	for k := 0; k < len(nodes); k++ {
		if nodes[k].id == id {
			return k
		}
	}
	return -1
}

/*
Esta función se encarga de buscar un nodo input en una interface de inputs.
*/
func findInput(input interface{}) string {
	return fmt.Sprintf("%v", input.(map[string]interface{})["connections"].([]interface{})[0].(map[string]interface{})["node"]) // fmt.Sprintf("%v", node1) permite convertir una interfaz en string
}

/*
Esta función se encarga de recibir un nodo y retornar su respuesta, dependiendo el tipo de nodo
si es un nodo Number, retorna el numero asociado al nodo, en caso de que sea un nodo Assign,
retorna la variable asociada al mismo, en caso de comparison o mathOperation se llama a la función
correspondiente.
*/
func typeNode(node Node, posNode int) string {
	nameNode := fmt.Sprintf("%v", node.name)
	if nameNode == "NodeNumber" || nameNode == "NodeAssign" {
		return fmt.Sprintf("%v", node.data.(map[string]interface{})["url"])
	} else if nameNode == "NodeComOp" {
		return comparison(posNode)
	} else if nameNode == "NodeMath" {
		return mathOperation("", posNode)
	}
	return ""
}

/*
Esta función se encarga de retornar el valor a ser asignado en la funcion assign.
*/
func valueAssigned(idNode int) string {
	return fmt.Sprintf("%v", nodes[idNode].data.(map[string]interface{})["url"])
}

/**
* Esta función se encarga de retornar una de cadena de 4 espacios, esto con el fin
* de hacer la identación, ya que el tabulador (\t) que trae por defecto go es de 8 espacios
* tener con 8 espacios no causa problemas en la ejecución, pero me parece mas ameno que tenga 4
 */
/*
func ident() string {
	return "	"
}*/

/*
Esta función se encarga de ordenar los nodos. Siendo C un conjunto de nodos,
C1, será el nodo que inicia las relación entre los nodos del conjunto y CN
el nodo que las finaliza.
Fanalmente retorna el conjunto de nodos ordenados.
*/
func sortNodes(nodeAux []Node) []Node {
	var sort []Node
	var posI int
	var isIf bool
	var father bool     // se usa para identificar si es un CN
	var brothers []Node // es un slice de Node, donde se guardaran los hermanos para cada iteración
	n := -1             // n será una variable acumuladora, por cada nodo que se agregue a sort(Node), esta aumentará en 1

	for {
		if len(nodeAux) == 1 { // Si queda un solo nodo, este se agrega y se rompe el ciclo
			sort = append(sort, nodeAux[0])
			break
		} else if isIf && nodeAux[posI].name == "NodeElse" || !isIf && nodeAux[posI].inputs.(map[string]interface{})["input_1"] == nil { // se busca el padre del arbol
			//if nodeAux[posI].name == "NodeAssign"
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
					thisOutput := findInput(thisNode)
					for i := 0; i < len(nodeAux); i++ { // se busca si ese padre tiene un hermano
						if nodeAux[i].outputs.(map[string]interface{})["output_1"] == nil || len(nodeAux[i].outputs.(map[string]interface{})["output_1"].(map[string]interface{})["connections"].([]interface{})) != 0 {

							if nodeAux[i].name != "NodePrint" && findInput(nodeAux[i].outputs.(map[string]interface{})["output_1"]) == thisOutput { // se verifica si la salida de X nodo es igual a la del nodo padre, si se cumple entonces son hermanos
								// Aqui ocurre un caso atipico con el Node assign, ya que existe la probabilidad de que su padre no haya sido agregado al slice.
								if nodeAux[i].inputs == nil { // si no tiene input, esto quiere decir que es un hermano que no tiene padre, por lo anterior es un CN
									brothers = append(brothers, nodeAux[i])
								} else { // si tiene padre, esto implica que no es un CN, por lo anterior se rompe el ciclo y se sigue buscando el CN
									brothers = nil
									father = false
									break
								}
								for j := 0; j < len(brothers); j++ {
									sort = append(sort, brothers[j])
									nodeAux = RemoveIndex(nodeAux, findNode(brothers[j].id, nodeAux))
									n++
								}
							}
						}
					}
					if !father { // como no es CN se rompe el otro ciclo
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
}

/*
Esta función se encarga de eliminar un elemento de una lista.
*/
func RemoveIndex(s []Node, index int) []Node {
	if len(s) == 1 {
		return nil
	}
	return append(s[:index], s[index+1:]...)
}

/*
{3 NodeComOp map[] map[output_1:map[connections:[map[node:1 output:input_2]]]] map[class:NodeComOp data:map[method:equals] html:NodeComOp id:3 inputs:map[] name:NodeComOp outputs:map[output_1:map[connections:[]]] pos_x:95 pos_y:318 typenode:vue]}
{12 NodeAssign map[input_1:map[connections:[map[input:output_1 node:8]]]] map[output_1:map[connections:[map[node:1 output:input_1]]]] map[url:y]}
{1 NodeIf map[input_1:map[connections:[map[input:output_1 node:12]]] input_2:map[connections:[map[input:output_1 node:3]]] input_3:map[connections:[map[input:output_1 node:4]]]] map[output_1:map[connections:[map[node:5 output:input_1]]] output_2:map[connections:[map[node:6 output:input_1]]]] map[url:1]}
{5 NodeAssign map[input_1:map[connections:[map[input:output_1 node:1]]]] map[output_1:map[connections:[]]] map[url:x]} {6 NodeElse map[input_1:map[connections:[map[input:output_2 node:1]]]] map[output_1:map[connections:[map[node:7 output:input_1]]]] map[url:2]} {7 NodeAssign map[input_1:map[connections:[map[input:output_1 node:6]]]] map[output_1:map[connections:[map[node:11 output:input_1]]]] map[url:x]} {11 NodePrint map[input_1:map[connections:[map[input:output_1 node:7]]]] map[] map[]} {9 NodeNumber map[] map[output_1:map[connections:[map[node:8 output:input_1]]]] map[url:2]} {10 NodeNumber map[] map[output_1:map[connections:[map[node:8 output:input_2]]]] map[url:1]} {8 NodeMath map[input_1:map[connections:[map[input:output_1 node:9]]] input_2:map[connections:[map[input:output_1 node:10]]]] map[output_1:map[connections:[map[node:12 output:input_1]]]] map[class:NodeMath data:map[method:add] html:NodeMath id:8 inputs:map[input_1:map[connections:[]] input_2:map[connections:[]]] name:NodeMath outputs:map[output_1:map[connections:[]]] pos_x:226 pos_y:118 typenode:vue url:3]} {4 NodeNumber map[] map[output_1:map[connections:[map[node:1 output:input_3]]]] map[url:3]}
*/
//{10 NodeNumber map[] map[output_1:map[connections:[map[node:8 output:input_2]]]] map[url:1]}
//{13 NodeMath map[input_1:map[connections:[map[input:output_1 node:14]]] input_2:map[connections:[map[input:output_1 node:15]]]] map[output_1:map[connections:[map[node:8 output:input_1]]]] map[class:NodeMath data:map[method:less] html:NodeMath id:13 inputs:map[input_1:map[connections:[]] input_2:map[connections:[]]] name:NodeMath outputs:map[output_1:map[connections:[]]] pos_x:53 pos_y:124 typenode:vue url:10]} {8 NodeMath map[input_1:map[connections:[map[input:output_1 node:13]]] input_2:map[connections:[map[input:output_1 node:10]]]] map[output_1:map[connections:[map[node:1 output:input_1]]]] map[class:NodeMath data:map[method:add] html:NodeMath id:8 inputs:map[input_1:map[connections:[]] input_2:map[connections:[]]] name:NodeMath outputs:map[output_1:map[connections:[]]] pos_x:226 pos_y:118 typenode:vue url:11]}
