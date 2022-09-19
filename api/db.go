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
	Languaje    string `json:"languaje,omitempty"`
}

type dataProgramUp struct {
	Uid  string
	Body string
}

func start_dgraph(option int, data_resp string) string {

	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	if option == 1 { // Mutation, se crea un nuevo programa

		// se convierte el idx(string) -> en una estructura Program
		data := Program{}
		json.Unmarshal([]byte(data_resp), &data)
		fmt.Println(data)

		op := &api.Operation{}
		op.Schema = `
			Name: string @index(exact) .
			ProgramName: string .
			Body: string .
			Languaje: string .
		`

		ctx := context.Background()
		err = dg.Alter(ctx, op)
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
		return "{success: Successfully added the new programa.}"

	}

	idx := string(data_resp)
	if option == 2 { // Query, obtiene un programa, lo filtra por su id
		// Assigned uids for nodes which were created would be returned in the assigned.Uids map.
		//variables := map[string]string{"$id1": assigned.Uids["alice"]}

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

	if option == 3 { // query, obtiene todos los programas almacenados en la base de datos, con la distinción de su nombre
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
		if err != nil {
			log.Fatal(err)
		}

		type Root struct {
			Fo []Program `json:"foo"`
		}

		var r Root
		err = json.Unmarshal(resp.Json, &r)
		if err != nil {
			log.Fatal(err)
		}

		return string(resp.Json)
	}

	if option == 4 { // Mutation, se edita un programa

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

		return "{success: Successfully updated the programa.}"
	}
	return ""
}
