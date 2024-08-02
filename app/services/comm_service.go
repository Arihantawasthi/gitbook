package services

import (
    "gitbook/app/types"
	"gitbook/utils"
	"strings"
)

type CommService struct {}

func NewCommService() CommService {
    return CommService{}
}

func (s *CommService) GetRepoCommits(repoPath, branch string) ([]types.Log, error) {
    commitLogList := []types.Log{}
    commitLog := types.Log{}
    output, err := utils.RunCommand("git", repoPath, "log", "--format=%H|%an|%ad|%s", "--date=format:%b %d, %Y %I:%M %p", branch)
    if err != nil {
        return nil, err
    }
    commits := strings.Split(output, "\n")
    for _, commit := range commits {
        commitParts := strings.Split(commit, "|")
        commitLog.Hash = commitParts[0]
        commitLog.Author = commitParts[1]
        commitLog.Timestamp = commitParts[2]
        commitLog.Message = commitParts[3]
        commitLogList = append(commitLogList, commitLog)
    }

    return commitLogList, nil
}
