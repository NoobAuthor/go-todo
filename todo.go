package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// TodoItem defines the structure for an API item
type TodoItem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// TodoStore holds the to-do list
var todoStore = make(map[string]TodoItem)
var mutex = &sync.Mutex{}

// GetTodoHandler handles the retrieval of to-do items
func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	var items []TodoItem
	mutex.Lock()
	for _, item := range todoStore {
		items = append(items, item)
	}
	mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// AddTodoHandler handles adding new to-do items
func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	var item TodoItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	todoStore[item.ID] = item
	mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// UpdateTodoHandler handles updates to existing to-do items
func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedItem TodoItem
	err := json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	item, exists := todoStore[id]
	if !exists {
		http.Error(w, "Todo item not found", http.StatusNotFound)
		mutex.Unlock()
		return
	}
	item.Description = updatedItem.Description
	item.Completed = updatedItem.Completed
	todoStore[id] = item
	mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// DeleteTodoHandler handles the deletion of to-do items
func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	mutex.Lock()
	_, exists := todoStore[id]
	if !exists {
		http.Error(w, "Todo item not found", http.StatusNotFound)
		mutex.Unlock()
		return
	}
	delete(todoStore, id)
	mutex.Unlock()
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todos", GetTodoHandler).Methods("GET")
	r.HandleFunc("/todos", AddTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id}", UpdateTodoHandler).Methods("PUT")
	r.HandleFunc("/todos/{id}", DeleteTodoHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
