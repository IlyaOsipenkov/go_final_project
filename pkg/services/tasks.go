package services

import (
	"database/sql"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTasks(db *sql.DB, search string, limit int) ([]Task, error) {
	var query string
	var args []any

	parsedDate, err := time.Parse("02.01.2006", search)
	if err == nil {
		query = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE CAST(date AS TEXT) = ?
		ORDER BY date ASC
		LIMIT ?`
		args = append(args, parsedDate.Format("20060102"), limit)
	} else if search != "" {
		query = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date ASC
		LIMIT ?`
		likeParam := "%" + search + "%"
		args = append(args, likeParam, likeParam, limit)
	} else {
		query = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?`
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

func AddTask(db *sql.DB, date, title, comment, repeat string) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, date, title, comment, repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
