package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"practice2/models"
	"practice2/storage"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		idStr := r.URL.Query().Get("id")

		if idStr == "" {
			var list []models.Task
			for _, t := range storage.Tasks {
				list = append(list, t)
			}
			json.NewEncoder(w).Encode(list)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, 400)
			return
		}

		task, ok := storage.Tasks[id]
		if !ok {
			http.Error(w, `{"error":"task not found"}`, 404)
			return
		}

		json.NewEncoder(w).Encode(task)

	case http.MethodPost:

		var body struct {
			Title string `json:"title"`
		}

		json.NewDecoder(r.Body).Decode(&body)

		if body.Title == "" {
			http.Error(w, `{"error":"invalid title"}`, 400)
			return
		}

		storage.Mutex.Lock()
		id := storage.IDCounter
		storage.IDCounter++

		task := models.Task{
			ID:    id,
			Title: body.Title,
			Done:  false,
		}

		storage.Tasks[id] = task
		storage.Mutex.Unlock()

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(task)

	case http.MethodPatch:

		idStr := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, 400)
			return
		}

		var body struct {
			Done bool `json:"done"`
		}

		json.NewDecoder(r.Body).Decode(&body)

		task, ok := storage.Tasks[id]
		if !ok {
			http.Error(w, `{"error":"task not found"}`, 404)
			return
		}

		task.Done = body.Done
		storage.Tasks[id] = task

		json.NewEncoder(w).Encode(map[string]bool{
			"updated": true,
		})
	}
}
