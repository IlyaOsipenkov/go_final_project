package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	init_db "github.com/ilyaosipenkov/practicum_final_project/pkg/db"
	handlers "github.com/ilyaosipenkov/practicum_final_project/pkg/handlers"
)

func main() {

	dbInstance := init_db.InitializeDB()
	fmt.Println("DB initialized succsessfully", dbInstance)

	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		filePath := "./web" + r.URL.Path

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Println("file not exist", filePath)
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, filePath)
	})

	r.Get("/api/nextdate", handlers.NextDateHandler)
	r.Post("/api/task", handlers.TaskHandler(dbInstance))
	r.Get("/api/task", handlers.TaskHandler(dbInstance))
	r.Put("/api/task", handlers.TaskHandler(dbInstance))
	r.Get("/api/tasks", handlers.TasksHandler(dbInstance))
	r.Post("/api/task/done", handlers.TaskDoneHandler(dbInstance))
	r.Delete("/api/task", handlers.DeleteTaskHandler(dbInstance))

	fmt.Println("port running on :7540")
	if err := http.ListenAndServe(":7540", r); err != nil {
		log.Fatal(err)
	}

}
