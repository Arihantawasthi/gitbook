package handler

import (
	"encoding/json"
	"gitbook/app/logic"
	"net/http"
)

func GetAllRepos(w http.ResponseWriter, r *http.Request) error {
	repoList, err := logic.GetRepoListFromDir()
	if err != nil {
		return err
	}

    response, err := logic.GetRepoDetails(repoList)
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(response)
	return nil
}
