package handler

import (
	"gitbook/app/services"
	"gitbook/app/types"
	"gitbook/utils"
	"net/http"
	"os"
)

type RepoHandler struct {
	repoPath string
	svc      services.RepoService
	logger   utils.Logger
}

func NewRepoHandler(logger utils.Logger) *RepoHandler {
	service := services.NewRepoService()
	return &RepoHandler{
		repoPath: os.Getenv("REPO_DIR"),
		svc:      service,
		logger:   logger,
	}
}

func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) error {
    h.logger.Info("incoming request", "handler: GetAllRepos", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	repoList, err := h.svc.GetRepoList(h.repoPath)
	if err != nil {
        h.logger.Info(err.Error(), "handler: GetRepoList", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("cannot read the git server directory", http.StatusServiceUnavailable)
	}

	response, err := h.svc.GetRepoDetails(h.repoPath, repoList)
	if err != nil {
        h.logger.Info(err.Error(), "handler: GetRepoDetails", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues", http.StatusServiceUnavailable)
	}

	jsonResponse := types.JsonResponse[[]types.RepoDetails]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the repositories",
		Data:          response,
	}
    h.logger.Info("request completed", "handler: GetAllRepos", r.Method, r.URL.Path, r.UserAgent(), "")
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}
