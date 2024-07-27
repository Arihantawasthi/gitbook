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
		repoDetails, err := getRepoDetails(h.repoPath, name)
		if err != nil {
			return utils.RaiseHTTPError(err, "cannot get commit history", 500)
		}
		createdAt := repoDetails["firstCommit"][0]
		author := repoDetails["firstCommit"][1]
		lastCommit := repoDetails["lastCommit"][0]
		name = strings.TrimSuffix(name, ".git")
		response = append(response, types.RepoResponse{Name: name, Desc: string(desc), CreatedAt: createdAt, Author: author, LastCommitAt: lastCommit})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}

func getRepoDetails(repoPath, name string) (map[string][]string, error) {
	repoDetails := make(map[string][]string)
	gitDir := fmt.Sprintf("--git-dir=%s/%s", repoPath, name)
	output, err := RunCommand("git", gitDir, "log", "--format=%ad|%an", "--date=format:%b %d, %Y %I:%M %p", "--reverse")
	fmt.Println(output)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	commits := strings.Split(output, "\n")
	repoDetails["firstCommit"] = strings.Split(commits[0], "|")
	repoDetails["lastCommit"] = strings.Split(commits[len(commits)-1], "|")
	return repoDetails, nil
}

func RunCommand(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	outputStr := strings.TrimSuffix(string(output), "\n")
	return outputStr, nil
}
