package handler

import (
	"NoobAuthor/go-todo/cmd/internal/app/model"
	"NoobAuthor/go-todo/cmd/internal/app/store"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	var items []model.TodoItem
	store.Mutex.Lock()
	for _, item := range store.TodoStore {
		items = append(items, item)
	}
	store.Mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	var item model.TodoItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store.Mutex.Lock()
	store.TodoStore[item.ID] = item
	store.Mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedItem model.TodoItem
	err := json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store.Mutex.Lock()
	item, exists := store.TodoStore[id]
	if !exists {
		http.Error(w, "Todo item not found", http.StatusNotFound)
		store.Mutex.Unlock()
		return
	}
	item.Description = updatedItem.Description
	item.Completed = updatedItem.Completed
	store.TodoStore[id] = item
	store.Mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	store.Mutex.Lock()
	_, exists := store.TodoStore[id]
	if !exists {
		http.Error(w, "Todo item not found", http.StatusNotFound)
		store.Mutex.Unlock()
		return
	}
	delete(store.TodoStore, id)
	store.Mutex.Unlock()
	w.WriteHeader(http.StatusOK)
}
