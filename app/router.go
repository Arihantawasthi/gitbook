package app

import (
	"gitbook/app/handler"
	"net/http"
    "gitbook/utils"
)

func RegisterRoutes(router *http.ServeMux) {
    repoHandler := handler.NewRepoHandler()

	router.HandleFunc("GET /api/v1/repos", utils.HandlerWrapper(repoHandler.GetAllRepos))
}
