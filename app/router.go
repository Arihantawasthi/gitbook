package app

import (
	"net/http"
    "gitbook/app/handler"
)

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /api/v1/repos", handler.GetAllRepos)
}
