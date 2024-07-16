package logic

import (
	"fmt"
	"gitbook/app/types"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func GetRepoDetails(repoList []fs.DirEntry) ([]types.RepoResponse, error) {
	var resp []types.RepoResponse
	repoDetails := &types.RepoResponse{}
	for _, repo := range repoList {
		repoName := repo.Name()
		if !strings.HasSuffix(repoName, ".git") {
			continue
		}
		repoPath := fmt.Sprintf("%s/%s", os.Getenv("REPO_DIR"), repoName)
		desc, err := getRepoDesc(repoPath)
		if err != nil {
			slog.Error(fmt.Sprintf("error while reading desctiption for %s %v", repoName, err))
			return nil, err
		}

		repoDetails.Name = repoName
		repoDetails.Desc = desc
		resp = append(resp, *repoDetails)
	}
	return resp, nil
}

func GetRepoListFromDir() ([]fs.DirEntry, error) {
	repos, err := os.ReadDir(os.Getenv("REPO_DIR"))
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading directory: %v", err))
		return nil, err
	}

	return repos, nil
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
