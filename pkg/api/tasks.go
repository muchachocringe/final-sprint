package api

import (
	"net/http"

	"todo-app/pkg/db"
)

const defaultTaskLimit = 50

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(defaultTaskLimit)
	if err != nil {
		writeError(w, "Ошибка получения задач")
		return
	}

	writeJSON(w, TasksResponse{Tasks: tasks})
}
