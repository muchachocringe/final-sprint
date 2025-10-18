package api

import (
	"encoding/json"
	"net/http"
	"time"

	"todo-app/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка декодирования JSON")
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи")
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
			writeError(w, "Неверный формат даты")
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
			writeError(w, err.Error())
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
		writeError(w, "Ошибка добавления задачи")
		return
	}

	writeJSON(w, map[string]interface{}{"id": id})
}
