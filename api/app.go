package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	fmt.Println("Holaaa")
	port := "10000"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	log.Printf("Starting up on http://localhost:%s", port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Api Visual Programming"))
	})

	r.Mount("/posts", postsResource{}.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}
