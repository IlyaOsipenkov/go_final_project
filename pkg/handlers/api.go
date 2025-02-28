package api

import (
	"fmt"
	"net/http"
	"time"

	services "github.com/ilyaosipenkov/practicum_final_project/pkg/services"
)

type Response struct {
	NextDate string `json:"next_date,omitempty"`
	Error    string `json:"error,omitempty"`
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

	fmt.Fprint(w, nextDate)

	// response := Response{}

	// if err != nil {
	// 	response.Error = err.Error()
	// } else {
	// 	response.NextDate = nextDate
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}
