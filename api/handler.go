package api

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

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

func renderJSON(w http.ResponseWriter, v any) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("error marshaling json: %s", err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

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

	resId := ResponseId{Id: id}

	renderJSON(w, resId)
}

func (server *TaskServer) GetAllTaskshandler(w http.ResponseWriter, r *http.Request) {
	var allTasks []taskstore.Task
	allTasks = server.store.GetAllTasks()

	renderJSON(w, allTasks)
}

func (server *TaskServer) DeleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	server.store.DeleteAllTasks()
	w.Write([]byte("Deleted all tasks successfuly"))
}

func (server *TaskServer) GetTaskById(w http.ResponseWriter, r *http.Request) {
	stringId := r.PathValue("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(w, "error converting string to int: id", http.StatusBadRequest)
		return
	}

	task, err := server.store.GetTask(id)
	if err != nil {
		http.Error(w, "error fetching task details", http.StatusBadRequest)
		return
	}

	renderJSON(w, task)
}

func (server *TaskServer) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	stringId := r.PathValue("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(w, "error converting string to int: id (delete task)", http.StatusBadRequest)
		return
	}

	err = server.store.DeleteTask(id)
	if err != nil {
		errMsg := fmt.Sprintf("error deleting task by id: %v", err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	successMsg := fmt.Sprintf("Deleted task with id: %d", id)
	w.Write([]byte(successMsg))
}

func (server *TaskServer) TagHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	task := server.store.GetTasksByTag(tag)

	renderJSON(w, task)
}

func (server *TaskServer) DueHandler(w http.ResponseWriter, r *http.Request) {
	year := r.PathValue("year")
	month := r.PathValue("month")
	day := r.PathValue("day")

	timeString := fmt.Sprintf("%v-%v-%v", year, month, day)
	time, err := time.Parse("2006-01-02", timeString)
	if err != nil {
		errMsg := fmt.Sprintf("error parsing time: %s", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	task := server.store.GetTasksByDueDate(time)

	renderJSON(w, task)
}
