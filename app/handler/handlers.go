package handler

import (
	"fmt"
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
		h.logger.Error(err.Error(), "repo_service: GetRepoList", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("cannot read the git server directory", http.StatusServiceUnavailable)
	}

	response, err := h.svc.GetRepoDetails(h.repoPath, repoList)
	if err != nil {
		h.logger.Error(err.Error(), "repo_service: GetRepoDetails", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: cannot fetch repo details", http.StatusServiceUnavailable)
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

func (h *RepoHandler) GetRepoObjects(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("incoming request", "handler: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	var repoObjects types.RepoObjects
	repoName := r.PathValue("name")
	objectType := r.PathValue("type")
    repoDir := fmt.Sprintf("--git-dir=%s/%s.git", h.repoPath, repoName)
	path := utils.ExtractRepoPath(r.URL.Path)

	objects, err := h.svc.GetRepoObjects(repoDir, r.PathValue("branch"), path)
	if err != nil {
		h.logger.Error(err.Error(), "repo_service: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: not able to read objects", http.StatusServiceUnavailable)
	}

	branches, err := h.svc.GetRepoBranches(h.repoPath, repoName)
	if err != nil {
		h.logger.Error(err.Error(), "repo_service: GetRepoBranches", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: not able to read branches", http.StatusServiceUnavailable)
	}

	desc, err := utils.RunCommand("cat", fmt.Sprintf("%s/%s.git/description", h.repoPath, repoName))
	if err != nil {
		h.logger.Error(err.Error(), "utils: RunCommand", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: not able to read description", http.StatusServiceUnavailable)
	}

    rawLines, err := h.svc.GetBlobRawLines(repoDir, r.PathValue("branch"), path, objectType)
    if err != nil {
        return utils.RaiseHTTPError("skill issues: not able to read blob", http.StatusServiceUnavailable)
    }

	repoObjects.Name = repoName
	repoObjects.Desc = desc
	repoObjects.Branches = branches
    repoObjects.Blob = rawLines
	repoObjects.Data = objects
	jsonResponse := types.JsonResponse[types.RepoObjects]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the repository objects",
		Data:          repoObjects,
	}
    h.logger.Info("request completed", "handler: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}
