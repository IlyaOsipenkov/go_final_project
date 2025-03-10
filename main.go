package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	init_db "github.com/ilyaosipenkov/practicum_final_project/pkg/db"
	handlers "github.com/ilyaosipenkov/practicum_final_project/pkg/handlers"
)

func main() {

	dbInstance := init_db.InitializeDB()
	fmt.Println("DB initialized succsessfully", dbInstance)

	r := chi.NewRouter()

	r.Route("/js", func(r chi.Router) {
		r.Get("/*", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))).ServeHTTP)

	})

	r.Route("/css", func(r chi.Router) {
		r.Get("/*", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))).ServeHTTP)

	})

	r.Get("/favicon.ico", http.FileServer(http.Dir("./web")).ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	r.Get("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/login.html")
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
