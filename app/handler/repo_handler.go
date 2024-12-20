package handler

import (
	"fmt"
	"gitbook/app/services"
	"gitbook/app/storage"
	"gitbook/app/types"
	"gitbook/utils"
	"net/http"
	"os"
    "strconv"
	"strings"
	"sync"
)

type RepoHandler struct {
	repoPath string
	svc      services.RepoService
	logger   utils.Logger
}

func NewRepoHandler(logger utils.Logger) *RepoHandler {
	service := services.NewRepoService()
	return &RepoHandler{
		repoPath: os.Getenv("REPO_DIR"),
		svc:      service,
		logger:   logger,
	}
}

// TODO: Add repositories in the database itself, since we are using pagination
func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) error {
    var limit int
	h.logger.Info("incoming request", "handler: GetAllRepos", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	repoList, err := h.svc.GetRepoList(h.repoPath)
	if err != nil {
		h.logger.Error(err.Error(), "repo_service: GetRepoList", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("cannot read the git server directory", http.StatusServiceUnavailable)
	}

    queryParam := r.URL.Query()["limit"]
    if len(queryParam) > 0 {
        limit, err = strconv.Atoi(queryParam[0])
        if err != nil {
            h.logger.Error(err.Error(), "repo_service: GetRepoList", r.Method, r.URL.Path, r.UserAgent(), r.Body)
            return utils.RaiseHTTPError("incorrect params", http.StatusBadRequest)
        }
        if limit == 0 {
            h.logger.Error("limit is passed as 0", "repo_service: GetRepoList", r.Method, r.URL.Path, r.UserAgent(), r.Body)
            return utils.RaiseHTTPError("incorrect params: limit can't be zero", http.StatusBadRequest)
        }
    }
    limit = min(limit, len(repoList))
    repoList = repoList[:limit]

	response, err := h.svc.GetRepoDetails(h.repoPath, repoList)
	if err != nil {
		h.logger.Error(err.Error(), "repo_service: GetRepoDetails", r.Method, r.URL.Path, r.UserAgent(), r.Body)
		return utils.RaiseHTTPError("skill issues: cannot fetch repo details", http.StatusServiceUnavailable)
	}

	jsonResponse := types.JsonResponse[[]types.RepoDetails]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the repositories",
		Data:          response,
	}
	h.logger.Info("request completed", "handler: GetAllRepos", r.Method, r.URL.Path, r.UserAgent(), "")
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}

func (h *RepoHandler) GetRepoObjects(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("incoming request", "handler: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	var (
		repoObjects types.RepoObjects
		wg          sync.WaitGroup
		errChan     = make(chan error, 3)
	)

	repoName := r.PathValue("name")
	objectType := r.PathValue("type")
	repoDir := fmt.Sprintf("--git-dir=%s/%s.git", h.repoPath, repoName)
	path := utils.ExtractRepoPath(r.URL.Path)
	if objectType == "blob" {
		path = strings.TrimSuffix(path, "/")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		objects, err := h.svc.GetRepoObjects(repoDir, r.PathValue("branch"), path)
		if err != nil {
			h.logger.Error(err.Error(), "repo_service: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
			errChan <- utils.RaiseHTTPError("skill issues: not able to read objects", http.StatusServiceUnavailable)
		}
		repoObjects.Data = objects
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		branches, err := h.svc.GetRepoBranches(h.repoPath, repoName)
		if err != nil {
			h.logger.Error(err.Error(), "repo_service: GetRepoBranches", r.Method, r.URL.Path, r.UserAgent(), r.Body)
			errChan <- utils.RaiseHTTPError("skill issues: not able to read branches", http.StatusServiceUnavailable)
		}
		repoObjects.Branches = branches
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		desc, err := utils.RunCommand("cat", fmt.Sprintf("%s/%s.git/description", h.repoPath, repoName))
		if err != nil {
			h.logger.Error(err.Error(), "utils: RunCommand", r.Method, r.URL.Path, r.UserAgent(), r.Body)
			errChan <- utils.RaiseHTTPError("skill issues: not able to read description", http.StatusServiceUnavailable)
		}
		repoObjects.Desc = desc
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rawLines, err := h.svc.GetBlobRawLines(repoDir, r.PathValue("branch"), path, objectType)
		if err != nil {
            logMsg := fmt.Sprintf("skill issues: %s", err.Error())
			h.logger.Error(logMsg, "utils: RunCommand", r.Method, r.URL.Path, r.UserAgent(), r.Body)
			errChan <- utils.RaiseHTTPError(err.Error(), http.StatusServiceUnavailable)
		}
		repoObjects.Blob = rawLines
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	repoObjects.Name = repoName
	jsonResponse := types.JsonResponse[types.RepoObjects]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the repository objects",
		Data:          repoObjects,
	}
	h.logger.Info("request completed", "handler: GetRepoObjects", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	return nil
}

func (h *RepoHandler) GetStats(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("incoming request", "handler: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
    stats, err := storage.GetStats()
    if err != nil {
		h.logger.Error(err.Error(), "utils: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
    }

	jsonReponse := types.JsonResponse[types.AggStats]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the stats",
		Data:          *stats,
	}
	h.logger.Info("request completed", "handler: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	utils.WriteJson(w, http.StatusOK, jsonReponse)
	return nil
}
