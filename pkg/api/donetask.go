package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, http.StatusInternalServerError, "db error: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid repeat: "+err.Error())
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, http.StatusInternalServerError, "db error: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
