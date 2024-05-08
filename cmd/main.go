package main

import (
	"NoobAuthor/go-todo/cmd/internal/app/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todos", handler.GetTodoHandler).Methods("GET")
	r.HandleFunc("/todos", handler.AddTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id}", handler.UpdateTodoHandler).Methods("PUT")
	r.HandleFunc("/todos/{id}", handler.DeleteTodoHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
