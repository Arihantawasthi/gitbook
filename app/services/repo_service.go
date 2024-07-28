package services

import (
	"fmt"
	"gitbook/app/types"
	"gitbook/utils"
	"strings"
)

type RepoService struct{}

func NewRepoService() RepoService{
    return RepoService{}
}

func (s *RepoService) GetRepoList(repoPath string) ([]string, error) {
	output, err := utils.RunCommand("ls", repoPath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(output, "\n")
	return lines, nil
}

func (s *RepoService) GetRepoDetails(repoPath string, repoList []string) ([]types.RepoResponse, error) {
	var repoDetails []types.RepoResponse
	for _, name := range repoList {
		descPath := fmt.Sprintf("%s/%s/description", repoPath, name)
		desc, err := utils.RunCommand("cat", descPath)
		if err != nil {
			return nil, err
		}
		commitSummary, err := getFirstCommitAndAuthor(repoPath, name)
		if err != nil {
			return nil, err
		}
		repoDetails = append(
			repoDetails,
			types.RepoResponse{
				Name:         strings.TrimSuffix(name, ".git"),
				Desc:         desc,
				Author:       commitSummary["author"],
				CreatedAt:    commitSummary["firstCommit"],
				LastCommitAt: commitSummary["lastCommit"],
			})
	}
	return repoDetails, nil
}

func getFirstCommitAndAuthor(repoPath, name string) (map[string]string, error) {
	repoDetails := make(map[string]string)
	gitDir := fmt.Sprintf("--git-dir=%s/%s", repoPath, name)
	output, err := utils.RunCommand("git", gitDir, "log", "--format=%ad|%an", "--date=format:%b %d, %Y %I:%M %p", "--reverse")
	if err != nil {
		return nil, err
	}
	commits := strings.Split(output, "\n")
	firstCommitInfo := strings.Split(commits[0], "|")
	repoDetails["firstCommit"], repoDetails["author"] = firstCommitInfo[0], firstCommitInfo[1]
	repoDetails["lastCommit"] = strings.Split(commits[len(commits)-1], "|")[0]
	return repoDetails, nil
}
