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

func (s *RepoService) GetRepoDetails(repoPath string, repoList []string) ([]types.RepoDetails, error) {
	var repoDetails []types.RepoDetails
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
        defaultBranch, err := getDefaultBranch(repoPath, name)
        if err != nil {
            return nil, err
        }
		repoDetails = append(
			repoDetails,
			types.RepoDetails{
				Name:          strings.TrimSuffix(name, ".git"),
				Desc:          desc,
                DefaultBranch: defaultBranch,
				Author:        commitSummary["author"],
				CreatedAt:     commitSummary["firstCommit"],
				LastCommitAt:  commitSummary["lastCommit"],
			},
        )
	}
	return repoDetails, nil
}

func (s *RepoService) GetRepoObjects(repoDir, branch, path string) ([]types.Objects, error) {
    objects := []types.Objects{}
    output, err := utils.RunCommand("git", repoDir, "ls-tree", "--format=%(objecttype)|%(path)|%(objectsize)", branch, path)
    if err != nil {
        return nil, err
    }
    objectList := strings.Split(output, "\n")
    for _, object := range objectList {
        objectSplit := strings.Split(object, "|")
        pathParts := strings.Split(objectSplit[1], "/")
        path := pathParts[len(pathParts)-1]
        objects = append(objects, types.Objects{
            Type: objectSplit[0],
            FullPath: objectSplit[1],
            Path: path,
            Size: objectSplit[2],
        })
    }

    return objects, nil
}

func (s *RepoService) GetRepoBranches(repoPath, repoName string) ([]string, error) {
    repoDir := fmt.Sprintf("--git-dir=%s/%s.git", repoPath, repoName)
    output, err := utils.RunCommand("git", repoDir, "for-each-ref", "--format=%(refname:short)", "refs/heads")
    if err != nil {
        return nil, err
    }
    branches := strings.Split(output, "\n")
    return branches, nil
}

func (s *RepoService) GetBlobRawLines(repoDir, branch, path, objectType string) ([]string, error) {
    rawLines := []string{}
    if objectType != "blob" {
        return rawLines, nil
    }
    output, err := utils.RunCommand("git", repoDir, "show", "master:"+path)
    if err != nil {
        return nil, err
    }
    rawLines = strings.Split(output, "\n")
    return rawLines, nil
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

func getDefaultBranch(repoPath, name string) (string, error) {
    ref, err := utils.RunCommand("cat", fmt.Sprintf("%s/%s/HEAD", repoPath, name))
    if err != nil {
        return "", err
    }

    prefix := "ref: refs/heads/"
    if strings.HasPrefix(ref, prefix) {
        return strings.TrimPrefix(ref, prefix), nil
    }

    return "", fmt.Errorf("invalid ref format")
}
