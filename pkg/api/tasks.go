package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.Tasks(tasksLimit, search)
	if err != nil {
		writeError(w, "db error: "+err.Error())
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}