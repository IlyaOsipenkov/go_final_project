package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	services "github.com/ilyaosipenkov/practicum_final_project/pkg/services"
)

type Response struct {
	NextDate string `json:"next_date,omitempty"`
	Error    string `json:"error,omitempty"`
}

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			handleAddTask(w, r, db)

		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id == "" {
				response := map[string]string{"error": "id omitted"}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			task, err := services.GetTaskById(db, id)
			if err != nil {
				response := map[string]string{"error": err.Error()}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTS-8")
			json.NewEncoder(w).Encode(task)

		case http.MethodPut:
			var task services.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, `{"error":"Issue of decoding JSON"}`, http.StatusBadRequest)
			}
			if task.ID == "" {
				response := map[string]string{"error": "id of task omitted"}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}
			if task.Title == "" {
				response := map[string]string{"error": "task title is empty"}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			} else if len(task.Title) > 100 {
				response := map[string]string{"error": "title too long"}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}
			if len(task.Comment) > 300 {
				response := map[string]string{"error": "comment too long"}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			_, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"invalid froamt of date"}`, http.StatusBadRequest)
				return
			}

			if !services.IsValidRepeat(task.Repeat) {
				response := map[string]string{"error": "invalid repeat"}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			err = services.UpdateTask(db, task)
			if err != nil {
				response := map[string]string{"error": err.Error()}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTS-8")
			json.NewEncoder(w).Encode(map[string]string{})
		default:
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		}
	}
}

func TasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"invalid method"}`, http.StatusMethodNotAllowed)
			return
		}

		search := r.URL.Query().Get("search")

		tasks, err := services.GetTasks(db, search, 30)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"error of getting tasks: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []services.Task{}
		}

		w.Header().Set("content-type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"method not exist"}`, http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "id omitted"}`, http.StatusBadRequest)
			return
		}

		err := services.DeleteTask(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func TaskDoneHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"Invalid method"}`, http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "id omitted"}`, http.StatusBadRequest)
			return
		}

		task, err := services.GetTaskById(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		if task.Repeat == "" {
			err = services.DeleteTask(db, id)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
				return
			}
		} else {
			now := time.Now()
			nextDate, err := services.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "wrong next date: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			err = services.UpdateTaskDate(db, id, nextDate)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func handleAddTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"error of deserialization JSON"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"title of task is empty"}`, http.StatusBadRequest)
		return
	}

	now := time.Now().Format("20060102")

	if task.Date == "" {
		task.Date = now
	}

	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"invalid date of task"}`, http.StatusBadRequest)
		return
	}

	if task.Repeat == "d 1" && task.Date == now {
		//nothing
	} else if task.Repeat != "" && parsedDate.Before(time.Now()) {

		nextDate, err := services.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else if parsedDate.Before(time.Now()) {
		task.Date = now
	}

	id, err := services.AddTask(db, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"error of writing in DB: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, `{"error": "missing query parameters"}`, http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, `{"error": "invalid 'now' date format"}`, http.StatusBadRequest)
		return
	}

	nextDate, err := services.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, nextDate)
}
