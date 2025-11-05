package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"net/http"
	"strconv"
	"strings"
	
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, map[string]string{"error": "id is required"})
			return
		}
		task, err := db.GetTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, task)

	case http.MethodPut:
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSON(w, map[string]string{"error": "invalid json"})
			return
		}
		if strings.TrimSpace(task.ID) == "" {
			writeJSON(w, map[string]string{"error": "id is required"})
			return
		}
		if _, err := strconv.Atoi(task.ID); err != nil {
			writeJSON(w, map[string]string{"error": "invalid id"})
			return
		}
		if strings.TrimSpace(task.Title) == "" {
			writeJSON(w, map[string]string{"error": "task title is required"})
			return
		}
		if strings.TrimSpace(task.Repeat) != "" {
			if err := validateRepeatFormat(task.Repeat); err != nil {
				writeJSON(w, map[string]string{"error": "invalid repeat format"})
				return
			}
		}
		if err := processTaskDates(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error(), "flag": "1"})
			return
		}
		if err := db.UpdateTask(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if strings.TrimSpace(id) == "" {
			writeJSON(w, map[string]string{"error": "id is required"})
			return
		}
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})

	default:
		writeJSON(w, map[string]string{"error": "method is not supported"})
	}
}

