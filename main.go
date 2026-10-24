package main

import (
	"restapi/http"
	"restapi/todo"
)

func main() {
	todoList := todo.NewList()

	//fmt.Println(todoList.ListTasks())

	handler := http.NewHTTPHandlers(todoList)

	server := http.NewHTTPServer(handler)
	if err := server.StartServer(); err != nil {
		panic(err)
	}
}

//TODO:
// - не кидать панику в handlerах
// - оформить обработку ошибки в w.Writer, а не fmt
