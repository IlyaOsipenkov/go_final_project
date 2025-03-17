package services

import (
	"database/sql"
	"fmt"
	"regexp"
)

func IsValidRepeat(repeat string) bool {
	if repeat == "" {
		return true
	}
	//Здравствуйте, можно ли сделать как-то кроме регулярки?
	validRepeat := regexp.MustCompile(`^(d \d{1,3}|y|w ([1-7],?)+|m (-?\d{1,2},?)+(\s([1-9]|1[0-2],?)*)?)?$`)
	return validRepeat.MatchString(repeat)
}

func GetTaskById(db *sql.DB, id string) (*Task, error) {
	query := `
	SELECT id, date, title, comment, repeat
	FROM scheduler
	WHERE id = ?`

	var task Task
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	} else if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(db *sql.DB, task Task) error {
	query := `
	UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Task not found")
	}
	return nil
}

func DeleteTask(db *sql.DB, id string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task was not found")
	}
	return nil
}

func UpdateTaskDate(db *sql.DB, id, newDate string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	_, err := db.Exec(query, newDate, id)
	return err
}
