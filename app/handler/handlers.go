package handler

import (
	"encoding/json"
	"fmt"
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
	output, err := RunCommand("ls", h.repoPath)
	if err != nil {
		return utils.RaiseHTTPError(err, "cannot read the directory", 500)
	}
	lines := strings.Split(output, "\n")

	for _, name := range lines {
        desc, err := RunCommand("cat", fmt.Sprintf("%s/%s/description", h.repoPath, name))
        if err != nil {
            return utils.RaiseHTTPError(err, "cannot read description", 500)
        }
        name = strings.TrimSuffix(name, ".git")
        response = append(response, types.RepoResponse{Name: name, Desc: string(desc)})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}


func RunCommand(cmdName string , args ...string) (string, error) {
    cmd := exec.Command(cmdName, args...)
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    outputStr := strings.TrimSuffix(string(output), "\n")
    return outputStr, nil
}
