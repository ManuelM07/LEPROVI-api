package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Program struct {
	Uid         string `json:"uid,omitempty"`
	Name        string `json:"name,omitempty"`
	ProgramName string `json:"program_name,omitempty"`
	Body        string `json:"body,omitempty"`
	Languaje    string `json:"languaje,omitempty"`
}

type dataProgramUp struct {
	Uid  string
	Body string
}

type CancelFunc func()

func getDgraphClient() (*dgo.Dgraph, CancelFunc) {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}

func createProgram(dataResp string) string {
	dg, cancel := getDgraphClient()
	defer cancel()

	// se convierte el data(string) -> en una estructura Program
	data := Program{}
	json.Unmarshal([]byte(dataResp), &data)

	op := &api.Operation{}
	op.Schema = `
		Name: string @index(exact) .
		ProgramName: string .
		Body: string .
		Languaje: string .
	`

	ctx := context.Background()
	err := dg.Alter(ctx, op)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb
	dg.NewTxn().Mutate(ctx, mu)
	return `{"success": "Successfully added the new programa."}`
}

func getProgram(idx string) string {
	dg, cancel := getDgraphClient()
	defer cancel()

	variables := map[string]string{"$idx": idx}
	q := `query Me($idx: string){
			me(func: uid($idx)) {
				uid
				name
				program_name
				body
				languaje
			}
		}`

	resp, err := dg.NewTxn().QueryWithVars(context.Background(), q, variables)
	if err != nil {
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

	return string(resp.Json)
}

func listPrograms() (string, error) {
	dg, cancel := getDgraphClient()
	defer cancel()

	q := `{
		foo(func: has(program_name)) {
		  uid
		  name
		  program_name
		  body
		  languaje
		}
	  }`

	resp, err := dg.NewTxn().Query(context.Background(), q)
	if err != nil { // mal manejo del error
		return "", err
	}

	type Root struct {
		Fo []Program `json:"foo"` // cambiar foo
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return "", err
	}

	return string(resp.Json), nil
}

func updateProgram(data_resp string) string {
	dg, cancel := getDgraphClient()
	defer cancel()

	data := dataProgramUp{}
	json.Unmarshal([]byte(data_resp), &data)
	fmt.Println((data))

	q := `query Me($idx: string){
			  v as var(func: uid($idx))
			}`

	mu := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(`uid(v) <body> %q .`, data.Body)),
	}

	req := &api.Request{
		Query:     q,
		Mutations: []*api.Mutation{mu},
		CommitNow: true,
		Vars:      map[string]string{"$idx": data.Uid},
	}

	if _, err := dg.NewTxn().Do(context.Background(), req); err != nil {
		log.Fatal(err)
	}

	return `{"success": "Successfully updated the programa."}`
}
