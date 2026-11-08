package api

import (
	"encoding/json"
	"net/http"

	"go_final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}