package api

import (
	"net/http"
	"time"

	"todo-app/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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

	now := time.Now()

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, "Ошибка удаления задачи")
			return
		}
	} else {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error())
			return
		}
		if err := db.UpdateDate(id, next); err != nil {
			writeError(w, "Ошибка обновления задачи")
			return
		}
	}

	writeJSON(w, map[string]interface{}{})
}
