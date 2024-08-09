package app

import (
	"gitbook/app/handler"
	"gitbook/utils"
	"net/http"
)

func RegisterRoutes(router *http.ServeMux) {
    appLogger := &utils.SlogLogger{}
    repoHandler := handler.NewRepoHandler(appLogger)
    commHandler := handler.NewCommitHandler(appLogger)

	router.HandleFunc("GET /api/v1/repos", utils.HandlerWrapper(repoHandler.GetAllRepos))
    router.HandleFunc("GET /api/v1/repo/{name}/{type}/metadata/{branch}/", utils.HandlerWrapper(repoHandler.GetRepoObjects))
    router.HandleFunc("GET /api/v1/repo/logs/{name}/{branch}", utils.HandlerWrapper(commHandler.GetCommitHistory))
    router.HandleFunc("GET /api/v1/repo/commit/{name}/{hash}", utils.HandlerWrapper(commHandler.GetCommitDetails))
}
