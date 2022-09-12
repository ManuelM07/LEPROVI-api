package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type postsResource struct{}

type ProgramX struct {
	Code         string
	Languaje     string
	VersionIndex string
}

type ProgramUnCode struct {
	Code     map[string]interface{}
	Languaje string
}

func (rs postsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)               // GET /posts - Read a list of posts.
	r.Post("/", rs.Create)            // POST /posts - Create a new post.
	r.Post("/program", rs.GetProgram) // POST /posts - Get program.
	r.Post("/run", rs.Run)            // POST /posts - Run program.

	r.Route("/{id}", func(r chi.Router) {
		r.Use(PostCtx)
		r.Get("/", rs.Get) // GET /posts/{id} - Read a single post by :id.
		//r.Put("/", rs.Update)    // PUT /posts/{id} - Update a single post by :id.
		//r.Delete("/", rs.Delete) // DELETE /posts/{id} - Delete a single post by :id.
		// post - run
	})

	return r
}

// Request Handler - GET /posts - leer y listar todos los programas.
func (rs postsResource) List(w http.ResponseWriter, r *http.Request) {
	resp := strings.NewReader(start_dgraph(3, ""))

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp); err != nil {
		return
	}
}

// Request Handler - POST /posts - Crear nuevo programa.
func (rs postsResource) Create(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	resp := strings.NewReader(start_dgraph(1, string(reqBody)))

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp); err != nil {
		return
	}

}

// Obtener programa
func (rs postsResource) GetProgram(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	dataResp := ProgramUnCode{}
	json.Unmarshal([]byte(string(reqBody)), &dataResp)

	//resp := strings.NewReader(mapJson(string(reqBody)))
	resp := strings.NewReader(mapJson(dataResp.Code, dataResp.Languaje))

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Run, ejecuta un nuevo programa haciendo uso de la api de jdoodle
func (rs postsResource) Run(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	dataResp := ProgramX{}
	json.Unmarshal([]byte(string(reqBody)), &dataResp)
	//respCode := fmt.Sprintf("%s", reqBody)
	//languaje := "nodejs" //python3
	//versionIndex := "4" //4

	data := map[string]interface{}{
		"clientId":     goDotEnvVariable("CLIENT_ID"),
		"clientSecret": goDotEnvVariable("CLIENT_SECRET"),
		"script":       dataResp.Code,
		"language":     dataResp.Languaje,
		"versionIndex": dataResp.VersionIndex,
	}

	jsonData, _ := json.Marshal(data)
	respData := strings.NewReader(string(jsonData))

	resp, err := http.Post("https://api.jdoodle.com/v1/execute", "application/json", respData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Crea un nuevo contexto en ctx, el cual asocia el valor del id
func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Request Handler - GET /posts/{id} - leer y mostrar un programa por :id.
func (rs postsResource) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string) //obtener id

	resp := strings.NewReader(start_dgraph(2, id))

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp); err != nil {
		return
	}
}
