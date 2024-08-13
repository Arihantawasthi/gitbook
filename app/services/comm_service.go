package services

import (
	"gitbook/app/types"
	"gitbook/utils"
	"strings"
)

type CommService struct{}

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

		commitStat, err := getShortStat(repoPath, commitLog.Hash)
		if err != nil {
			return commitLogList, err
		}
		commitLog.FilesChanged = commitStat.FilesChanged
		commitLog.Deletions = commitStat.Deletions
		commitLog.Insertions = commitStat.Insertions
		commitLogList = append(commitLogList, commitLog)
	}

	return commitLogList, nil
}

func (s *CommService) GetFilesChangedInCommit(repoPath, hash string) ([]string, error) {
	output, err := utils.RunCommand("git", repoPath, "show", "--name-only", "--pretty=", hash)
	if err != nil {
		return nil, err
	}
	outputList := strings.Split(output, "\n")
	return outputList, nil
}

func (s *CommService) GetFilesDiff(repoPath, hash string, files []string) ([]types.DiffResponse, error) {
    filesDiff := []types.DiffResponse{}
	for _, file := range files {
		output, err := utils.RunCommand("git", repoPath, "show", "--pretty=", hash, "--", file)
		if err != nil {
			return nil, err
		}
        lines := strings.Split(output, "\n")
        diffRes := types.DiffResponse{
            FilePath: file,
            CodeLines: lines,
        }
        filesDiff = append(filesDiff, diffRes)
	}
	return filesDiff, nil
}

func getShortStat(repoPath, commitHash string) (types.LogStat, error) {
	logStat := types.LogStat{}
	output, err := utils.RunCommand("git", repoPath, "show", "--pretty=", "--shortstat", commitHash)
	if err != nil {
		return logStat, err
	}
	logStatParts := strings.Split(output, ",")
	logStat = extractLogStat(logStatParts)
	formattedLogStat := formatLogStat(logStat)
	return formattedLogStat, nil
}

func extractLogStat(logStatParts []string) types.LogStat {
	logStat := types.LogStat{}
	logStat.FilesChanged = logStatParts[0]
	logStat.Deletions = "0"
	logStat.Insertions = "0"
	if strings.Contains(logStatParts[1], "-") {
		logStat.Deletions = logStatParts[1]
	}
	if strings.Contains(logStatParts[1], "+") {
		logStat.Insertions = logStatParts[1]
	}
	if len(logStatParts) <= 2 {
		return logStat
	}
	if strings.Contains(logStatParts[2], "-") {
		logStat.Deletions = logStatParts[2]
	}
	if strings.Contains(logStatParts[2], "+") {
		logStat.Insertions = logStatParts[2]
	}

	return logStat
}

func formatLogStat(logStat types.LogStat) types.LogStat {
	filesChanged := strings.TrimPrefix(logStat.FilesChanged, " ")
	deletions := strings.TrimPrefix(logStat.Deletions, " ")
	insertions := strings.TrimPrefix(logStat.Insertions, " ")

	filesChangedParts := strings.Split(filesChanged, " ")
	insertionsParts := strings.Split(deletions, " ")
	deletionsParts := strings.Split(insertions, " ")

	logStat.FilesChanged = filesChangedParts[0]
	logStat.Insertions = insertionsParts[0]
	logStat.Deletions = deletionsParts[0]

	return logStat
}
