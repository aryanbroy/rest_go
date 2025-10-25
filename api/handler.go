package api

import (
	"net/http"

	"github.com/aryanbroy/rest_go/internal/taskstore"
)

type TaskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *TaskServer {
	store := taskstore.New()
	return &TaskServer{store: store}
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello there niga"))

}

func (server *TaskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the create task handler bitch"))
}
