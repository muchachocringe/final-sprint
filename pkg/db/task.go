package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.Get(&task, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача не найдена")
	}
	return &task, err
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(id string, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := db.Exec(query, date, id)
	return err
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= date('now') ORDER BY date LIMIT ?`
	err := db.Select(&tasks, query, limit)
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, err
}
