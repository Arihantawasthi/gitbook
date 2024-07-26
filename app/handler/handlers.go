package handler

import (
	"encoding/json"
	"gitbook/utils"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RepoHandler struct {
    repoPath string
}

func NewRepoHandler() *RepoHandler {
    return &RepoHandler{
        repoPath: os.Getenv("REPO_DIR"),
    }
}

type RepoResponse struct {
    RepoName string
    Desc string
    Author string
    CreatedAt time.Time
    LastCommit time.Time
}

func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) error {
    var response []RepoResponse
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
        response = append(response, RepoResponse{RepoName: name})
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    return nil
}
