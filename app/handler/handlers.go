package handler

import (
	"gitbook/app/services"
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
		return utils.RaiseHTTPError(err, "cannot read the git server directory", 500)
	}

	response, err := h.svc.GetRepoDetails(h.repoPath, repoList)
	if err != nil {
		return utils.RaiseHTTPError(err, "skill issues", 500)
	}

	utils.WriteJson(w, http.StatusOK, response)
	return nil
}
