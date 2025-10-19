package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"todo-app/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	afterNow := func(date, now time.Time) bool {
		return date.Format(DateFormat) > now.Format(DateFormat)
	}

	now := time.Now()
	today := now.Format(DateFormat)

	if task.Date == "" {
		task.Date = today
	}

	var t time.Time
	var err error
	if task.Date != today {
		t, err = time.Parse(DateFormat, task.Date)
		if err != nil {
			writeError(w, "Неверный формат даты", http.StatusBadRequest)
			return
		}
	} else {
		task.Date = today
		t, _ = time.Parse(DateFormat, today)
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			if strings.Contains(err.Error(), "Ошибка вычисления") {
				writeError(w, err.Error(), http.StatusInternalServerError)
			} else {
				writeError(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка добавления задачи", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"id": id})
}
