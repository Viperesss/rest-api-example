package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"restapi/todo"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
	}
}

/*
pattern: /tasks
methood: POST
info: json body

succeed:

	status code: 201
	return: json with new task

failed:

	status: 400, 409, 500...
	return: json with error + time
*/
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var TaskDTO TaskDTO

	// Читаем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&TaskDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	// Валидируем данные
	if err := TaskDTO.ValidateForCreate(); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	// Добавляем задау в общий лист
	todoTask := todo.NewTask(TaskDTO.Title, TaskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return

	}

	// Даем ответ
	b, err := json.MarshalIndent(todoTask, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
methood: GET
info: pattern

succeed:

	status code: 200
	return: json with task

failed:

	status: 400, 404, 500...
	return: json with error + time
*/
func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	task, err := h.todoList.GetTask(title)
	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return

	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks
methood: GET
info: -

succeed:

	status code: 200
	return: json with tasks

failed:

	status: 400, 500...
	return: json with error + time
*/
func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()

	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks?completed=false
methood: GET
info: query

succeed:

	status code: 200
	return: json with uncopmleted tasks

failed:

	status: 400, 500...
	return: json with error + time
*/
func (h *HTTPHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncomplitedTasks := h.todoList.ListUncompletedTasks()

	b, err := json.MarshalIndent(uncomplitedTasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
methood: PATCH
info: pattern + json body

succeed:

	status code: 200
	return: json with updated tasks

failed:

	status: 400, 409, 500...
	return: json with error + time
*/
func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	var (
		changedTask todo.Task
		err         error
	)

	if completeDTO.Complete {
		changedTask, err = h.todoList.CompleteTask(title)
	} else {
		changedTask, err = h.todoList.UncompleteTask(title)
	}

	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(changedTask, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
methood: DELETE
info: pattern

succeed:

	status code: 204
	return: -

failed:

	status: 400, 404, 500...
	return: json with error + time
*/
func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
