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
		//fmt.Println(node)
	}
	fmt.Println(nodes)

	// para acceder a la clave de un map de varios niveles
	//fmt.Println(m["1"].(map[string]interface{})["data"].(map[string]interface{})["url"])

}
