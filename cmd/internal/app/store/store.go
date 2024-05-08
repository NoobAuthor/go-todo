package store

import (
	"NoobAuthor/go-todo/cmd/internal/app/model"
	"sync"
)

var TodoStore = make(map[string]model.TodoItem)
var Mutex = &sync.Mutex{}
