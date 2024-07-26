package services

import (
    "gitbook/app/types"
)

type RepoService struct {
	repoPath string
}


func NewRepoService(repoPath string) *RepoService {
    return &RepoService{repoPath: repoPath}
}


func (r *RepoService) GetRepoListFromDir() {
    var response []types.RepoResponse
}

