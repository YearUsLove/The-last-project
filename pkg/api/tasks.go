package api

import (
	"net/http"
	"strings"

	"final_project/pkg/db"
)

const defaultLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{"error": "method is not supported"})
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	tasks, err := db.Tasks(defaultLimit, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
