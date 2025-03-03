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
		default:
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	}
}

func handleAddTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Error of desyrization JSON"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Title of task is empty"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"Invalid date of task"}`, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" {
		nextDate, err := services.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else if parsedDate.Before(now) {
		task.Date = now.Format("20060102")
	}

	parsedDate, _ = time.Parse("20060102", task.Date)
	if parsedDate.After(now) {
		task.Date = now.Format("20060102")
	}

	id, err := services.AddTask(db, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Error of writing in BD: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	// fmt.Fprint(w, id) wihtout JSON, so for tests
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid 'now' date format", http.StatusBadRequest)
		return
	}

	nextDate, err := services.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, nextDate) //wihtout JSON, so for tests

	//JSON format for TODOS complete

	// response := Response{}

	// if err != nil {
	// 	response.Error = err.Error()
	// } else {
	// 	response.NextDate = nextDate
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}
