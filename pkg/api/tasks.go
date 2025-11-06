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
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method is not supported"})
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	tasks, err := db.Tasks(defaultLimit, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}