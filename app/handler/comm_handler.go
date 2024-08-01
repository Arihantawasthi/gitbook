package handler

import (
	"gitbook/utils"
	"net/http"
	"os"
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


func (h *CommitHandler) GetCommitHistory(w http.ResponseWriter, r *http.Request) error {
    return nil
}
