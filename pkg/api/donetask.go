package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена")
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, "db error: "+err.Error())
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, "invalid repeat: "+err.Error())
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, "db error: "+err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}