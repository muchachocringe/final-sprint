package api

import (
	"net/http"

	"todo-app/pkg/db"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, "Ошибка получения задач")
		return
	}

	writeJSON(w, TasksResponse{Tasks: tasks})
}
