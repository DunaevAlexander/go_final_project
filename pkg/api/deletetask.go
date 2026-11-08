package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}