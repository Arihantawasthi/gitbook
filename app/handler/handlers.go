package handler

import (
	"encoding/json"
	"net/http"
)


func GetAllRepos(w http.ResponseWriter, r *http.Request) {
    items := map[string]int{"REPO": 1}
    json.NewEncoder(w).Encode(items)
}
