package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
)

type Program struct {
	Uid         string `json:"uid,omitempty"`
	Name        string `json:"name,omitempty"`
	ProgramName string `json:"program_name,omitempty"`
	Body        string `json:"body,omitempty"`
}

func start_dgraph(option int, idx string) string {

	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	if option == 1 { // Mutation
		p := Program{
			Uid:         "_:veronica",
			Name:        "Veronica",
			ProgramName: "Add",
			Body:        "def sum(a, b):\n\treturn a + b",
		}

		op := &api.Operation{}
		op.Schema = `
			name: string @index(exact) .
			ProgramName: string .
			Body: string .
		`

		ctx := context.Background()
		err = dg.Alter(ctx, op)
		if err != nil {
			log.Fatal(err)
		}

		mu := &api.Mutation{
			CommitNow: true,
		}
		pb, err := json.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}

		mu.SetJson = pb
		dg.NewTxn().Mutate(ctx, mu)
	}

	if option == 2 { // Query
		// Assigned uids for nodes which were created would be returned in the assigned.Uids map.
		//variables := map[string]string{"$id1": assigned.Uids["alice"]}
		variables := map[string]string{"$idx": idx}
		q := `query Me($idx: string){
			me(func: uid($idx)) {
			  	uid
			  	name
				program_name
				body
			}
		}`

		resp, err := dg.NewTxn().QueryWithVars(context.Background(), q, variables) //QueryWithVars, dg.NewTxn().Query(context.Background(), q)
		if err != nil {
			fmt.Println("Holaaa")
			log.Fatal(err)
		}

		type Root struct {
			Me []Program `json:"me"`
		}

		var r Root
		err = json.Unmarshal(resp.Json, &r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(resp.Json))
		return string(resp.Json)
	}
	return ""
}
