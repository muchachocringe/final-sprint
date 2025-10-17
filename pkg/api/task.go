package api

import (
	"encoding/json"
	"net/http"
	"time"

	"todo-app/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена")
		return
	}

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка декодирования JSON")
		return
	}

	if task.ID == "" {
		writeError(w, "Не указан идентификатор задачи")
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи")
		return
	}

	if task.Date != "" {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			writeError(w, "Неверный формат даты")
			return
		}

		now := time.Now().Format("20060102")
		if task.Date < now {
			writeError(w, "Дата не может быть в прошлом")
			return
		}
	}

	if task.Repeat != "" {
		now := time.Now()
		testDate := now.Format("20060102")
		if task.Date != "" {
			testDate = task.Date
		}
		if _, err := NextDate(now, testDate, task.Repeat); err != nil {
			writeError(w, err.Error())
			return
		}
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, "Ошибка обновления задачи")
		return
	}

	writeJSON(w, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор")
		return
	}

	_, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, "Ошибка удаления задачи")
		return
	}

	writeJSON(w, map[string]interface{}{})
}
