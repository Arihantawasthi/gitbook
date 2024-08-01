package handler

import (
	"fmt"
	"gitbook/app/types"
	"gitbook/utils"
	"net/http"
	"os"
	"strings"
)

type CommitHandler struct {
	repoPath string
	logger   utils.Logger
}

func NewCommitHandler(logger utils.Logger) *CommitHandler {
	return &CommitHandler{
		repoPath: os.Getenv("REPO_DIR"),
		logger:   logger,
	}
}

type Log struct {
	Hash      string `json:"commit_hash"`
	Author    string `json:"commit_author"`
	Message   string `json:"commit_message"`
	Timestamp string `json:"commit_timestamp"`
}

func (h *CommitHandler) GetCommitHistory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("incoming request", "handler: GetCommitHistory", r.Method, r.URL.Path, r.UserAgent(), r.Body)
    var commitLog Log
    commitLogs := []Log{}
	repoName := r.PathValue("name") + ".git"
	gitDir := fmt.Sprintf("--git-dir=%s/%s", h.repoPath, repoName)
	output, err := utils.RunCommand("git", gitDir, "log", "--format=%H|%an|%ad|%s", "--date=format:%b %d, %Y %I:%M %p", "master")
	if err != nil {
		h.logger.Error(err.Error(), "handler: GetCommitHistory", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: not able to get the logs", http.StatusServiceUnavailable)
	}
	logs := strings.Split(output, "\n")
	for _, log := range logs {
		logParts := strings.Split(log, "|")
        commitLog.Hash = logParts[0]
        commitLog.Author = logParts[1]
        commitLog.Message = logParts[2]
        commitLog.Timestamp = logParts[3]
        commitLogs = append(commitLogs, commitLog)
	}

    jsonResponse := types.JsonResponse[[]Log]{
        RequestStatus: 1,
        StatusCode: http.StatusOK,
        Msg: "Successfully retreived the repo logs",
        Data: commitLogs,
    }
    h.logger.Info("request completed", "handler: GetCommitHistory", r.Method, r.URL.Path, r.UserAgent(), r.Body)
    utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}
