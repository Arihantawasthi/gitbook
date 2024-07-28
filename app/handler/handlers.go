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
}

func NewRepoHandler() *RepoHandler {
	service := services.NewRepoService()
	return &RepoHandler{
		repoPath: os.Getenv("REPO_DIR"),
		svc:      service,
	}
}

func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) error {
	repoList, err := h.svc.GetRepoList(h.repoPath)
	if err != nil {
		return utils.RaiseHTTPError("cannot read the git server directory", http.StatusServiceUnavailable)
	}

	response, err := h.svc.GetRepoDetails(h.repoPath, repoList)
	if err != nil {
		return utils.RaiseHTTPError("skill issues", http.StatusServiceUnavailable)
	}

	jsonResponse := types.JsonResponse[[]types.RepoDetails]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the repositories",
		Data:          response,
	}
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}
