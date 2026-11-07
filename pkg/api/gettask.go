package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, task)
}