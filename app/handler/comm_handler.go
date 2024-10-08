package handler

import (
	"fmt"
	"gitbook/app/services"
	"gitbook/app/types"
	"gitbook/utils"
	"net/http"
	"os"
)

type CommitHandler struct {
	repoPath string
	svc      services.CommService
	logger   utils.Logger
}

func NewCommitHandler(logger utils.Logger) *CommitHandler {
	service := services.NewCommService()
	return &CommitHandler{
		repoPath: os.Getenv("REPO_DIR"),
		svc:      service,
		logger:   logger,
	}
}

func (h *CommitHandler) GetCommitHistory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("incoming request", "handler: GetCommitHistory", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	repoName := r.PathValue("name") + ".git"
	gitDir := fmt.Sprintf("--git-dir=%s/%s", h.repoPath, repoName)
	descPath := fmt.Sprintf("%s/%s/description", h.repoPath, repoName)
	desc, err := utils.RunCommand("cat", descPath)
	logs, err := h.svc.GetRepoCommits(gitDir, r.PathValue("branch"))
	if err != nil {
		return utils.RaiseHTTPError("skill issues: error in reading logs", http.StatusServiceUnavailable)
	}

	jsonResponse := types.JsonResponse[types.CommitHistory]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retreived the repo logs",
		Data: types.CommitHistory{
			Name: repoName,
			Desc: desc,
			Logs: logs,
		},
	}
	h.logger.Info("request completed", "handler: GetCommitHistory", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}

func (h *CommitHandler) GetFileCommits(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("incoming request", "handler: GetFilesCommits", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	repoName := r.PathValue("name") + ".git"
	gitDir := fmt.Sprintf("--git-dir=%s/%s", h.repoPath, repoName)
	file := r.PathValue("file")
	commits, err := h.svc.GetFileCommits(gitDir, file)
	if err != nil {
		return utils.RaiseHTTPError("skill issues: error in reading file commit history", http.StatusServiceUnavailable)
	}

	jsonResponse := types.JsonResponse[[]types.Log]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the file history",
		Data:          commits,
	}
	h.logger.Info("request completed", "handler: GetFileCommits", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}

func (h *CommitHandler) GetCommitDetails(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("incoming request", "handler: GetCommitDetails", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	repoName := r.PathValue("name") + ".git"
	gitDir := fmt.Sprintf("--git-dir=%s/%s", h.repoPath, repoName)
	filesChanged, err := h.svc.GetFilesChangedInCommit(gitDir, r.PathValue("hash"))
	if err != nil {
		return utils.RaiseHTTPError("skill issues: error in reading logs", http.StatusServiceUnavailable)
	}
	diff, err := h.svc.GetFilesDiff(gitDir, r.PathValue("hash"), filesChanged)
	jsonResponse := types.JsonResponse[[]types.DiffResponse]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the repo logs",
		Data:          diff,
	}
	h.logger.Info("request completed", "handler: GetCommitDetails", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}
