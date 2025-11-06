package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPut:
		handleEditTask(w, r)
	case http.MethodDelete:
		handleDeleteTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{"error": "method is not supported"})
	}
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, task)
}

func handleEditTask(w http.ResponseWriter, r *http.Request) {
	var incoming map[string]any
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid json"})
		return
	}

	task := db.Task{}
	if v, ok := incoming["id"]; ok {
		task.ID = fmt.Sprint(v)
	}
	if strings.TrimSpace(task.ID) == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	if _, err := strconv.Atoi(task.ID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid id"})
		return
	}
	if v, ok := incoming["date"]; ok {
		task.Date = fmt.Sprint(v)
	}
	if v, ok := incoming["title"]; ok {
		task.Title = fmt.Sprint(v)
	}
	if strings.TrimSpace(task.Title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task title is required"})
		return
	}
	if v, ok := incoming["comment"]; ok {
		task.Comment = fmt.Sprint(v)
	}
	if v, ok := incoming["repeat"]; ok {
		task.Repeat = fmt.Sprint(v)
	}
	if task.Repeat != "" {
		if err := validateRepeatFormat(task.Repeat); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, map[string]string{"error": "invalid repeat format"})
			return
		}
	}
	if err := processTaskDates(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error(), "flag": "1"})
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{})
}
