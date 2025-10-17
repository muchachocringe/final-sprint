package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"runtime"
)

func Init() {
	_, filename, _, _ := runtime.Caller(0)
	webDir := filepath.Join(filepath.Dir(filename), "../../web")

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string) {
	writeJSON(w, map[string]string{"error": message})
}
