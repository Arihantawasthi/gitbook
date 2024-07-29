package handler

import (
	"fmt"
	"gitbook/app/services"
	"gitbook/app/types"
	"gitbook/utils"
	"net/http"
	"os"
	"strings"
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
		h.logger.Error(err.Error(), "handler: GetRepoList", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("cannot read the git server directory", http.StatusServiceUnavailable)
	}

	response, err := h.svc.GetRepoDetails(h.repoPath, repoList)
	if err != nil {
		h.logger.Error(err.Error(), "handler: GetRepoDetails", r.Method, r.URL.Path, r.UserAgent(), r.Body)
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
	objects := []types.Objects{}
	repoObjects := types.RepoObjects{}
	repoName := r.PathValue("name") + ".git"
	repoDir := fmt.Sprintf("--git-dir=%s/%s", h.repoPath, repoName)
	output, err := utils.RunCommand("git", repoDir, "ls-tree", "--format=%(objecttype)|%(path)", "master")
	if err != nil {
        h.logger.Error(err.Error(), "handler: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: not able to repository objects", http.StatusServiceUnavailable)
	}
	desc, err := utils.RunCommand("cat", fmt.Sprintf("%s/%s/description", h.repoPath, repoName))
	if err != nil {
        h.logger.Error(err.Error(), "handler: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: not able to read description", http.StatusServiceUnavailable)
	}
	branchOutput, err := utils.RunCommand("git", repoDir, "for-each-ref", "--format=%(refname:short)", "refs/heads")
	branchList := strings.Split(branchOutput, "\n")
	if err != nil {
		return utils.RaiseHTTPError("skill issues: not able to get all the branches", http.StatusServiceUnavailable)
	}

	objectStr := strings.Split(output, "\n")
	for _, object := range objectStr {
		objectSplit := strings.Split(object, "|")
		objects = append(objects, types.Objects{Type: objectSplit[0], Path: objectSplit[1]})
	}
	repoObjects.Name = strings.TrimSuffix(repoName, ".git")
	repoObjects.Desc = desc
	repoObjects.Branches = branchList
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
