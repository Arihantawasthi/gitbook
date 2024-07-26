package handler

import (
	"encoding/json"
	"gitbook/app/types"
	"gitbook/utils"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type RepoHandler struct {
	repoPath string
}

func NewRepoHandler() *RepoHandler {
	return &RepoHandler{
		repoPath: os.Getenv("REPO_DIR"),
	}
}

func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) error {
	var response []types.RepoResponse
	cmd := exec.Command("ls")
	cmd.Dir = h.repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return utils.RaiseHTTPError(err, "cannot read the directory", 500)
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, name := range lines {
		name = strings.TrimSuffix(name, ".git")
        response = append(response, types.RepoResponse{Name: name})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}
