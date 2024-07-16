package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Repo struct {
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Author       string `json:"author"`
	CreatedAt    string `json:"created_at"`
	LastCommitAt string `json:"last_commit_at"`
}

func GetAllRepos(w http.ResponseWriter, r *http.Request) error {
    var response []Repo
	repos, err := os.ReadDir(os.Getenv("REPO_DIR"))
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading directory: %v", err))
		return err
	}

	for _, repo := range repos {
		repoName := repo.Name()
		if !strings.HasSuffix(repoName, ".git") {
			continue
		}
		repoPath := fmt.Sprintf("%s/%s", os.Getenv("REPO_DIR"), repoName)
		desc, err := getRepoDesc(repoPath)
		if err != nil {
			return err
		}
		// Get Repo Author
		// Get First Commit
		// Get Last Commit
        repoDetails := &Repo{
            Name: repoName,
            Desc: desc,
        }
        response = append(response, *repoDetails)
	}
    json.NewEncoder(w).Encode(response)
	return nil
}

func getRepoDesc(repoPath string) (string, error) {
	descPath := fmt.Sprintf("%s/%s", repoPath, "description")
	cmd := exec.Command("cat", descPath)
	output, err := cmd.Output()
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading description: %v", err))
		return "", err
	}

	return string(output), nil
}
