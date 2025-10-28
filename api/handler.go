package api

import (
	"encoding/json"
	"log"
	"mime"
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

// func TestHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello there niga"))
//
// }

func (server *TaskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task taskstore.Task

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(string(contentType))
	if err != nil {
		log.Println(err)
		http.Error(w, "error parsing content type", 400)
		return
	}

	if mediaType != "application/json" {
		http.Error(w, "request content type doesn't match the required content type", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&task)
	if err != nil {
		http.Error(w, "error decoding request body containing task", http.StatusBadRequest)
		return
	}

	id := server.store.CreateTask(task.Text, task.Tags, task.Due)
	type ResponseId struct {
		Id int `json:"id"`
	}
	jsonData, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, "error marshaling json data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(jsonData)
}

func (server *TaskServer) GetAllTaskskhandler(w http.ResponseWriter, r *http.Request) {
	var allTasks []taskstore.Task
	allTasks = server.store.GetAllTasks()

	jsonData, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, "error decoding all tasks to json object", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
