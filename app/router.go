package app

import (
	"gitbook/app/handler"
	"gitbook/utils"
	"net/http"
)

func RegisterRoutes(router *http.ServeMux) {
    appLogger := &utils.SlogLogger{}
    repoHandler := handler.NewRepoHandler(appLogger)

	router.HandleFunc("GET /api/v1/repos", utils.HandlerWrapper(repoHandler.GetAllRepos))
    router.HandleFunc("GET /api/v1/repos/{name}/{branch}/{path}", utils.HandlerWrapper(repoHandler.GetRepoObjects))
}
