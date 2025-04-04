package handler

import (
	"database/sql"
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
	"time"
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
	h.logger.Info("incoming request", "handler: GetAllRepos", r.Method, r.URL.Path, r.UserAgent(), r.Body)
    limitStr := r.URL.Query().Get("limit")
    pageStr := r.URL.Query().Get("page")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    offset := (page - 1) * limit

    repos, err := storage.GetRepos(limit, offset)
    if err != nil {
        h.logger.Error(err.Error(), "utils: GetRepos", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        return err
    }

    jsonResponse := types.JsonResponse[[]types.RepoDetails]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the repos",
		Data:          repos,
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
    today := time.Now().UTC()
    todayDate := today.Truncate(24 * time.Hour)

    comparisonDate := todayDate.AddDate(0, 0, -7)
    currentStats, err := storage.GetStatsForADate(todayDate)
    if err != nil {
		h.logger.Error(err.Error(), "current stats, utils: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        if err != sql.ErrNoRows {
            return err
        }

	    h.logger.Info("calling GetLatestStats", "handler: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        currentStats, err = storage.GetLatestStats()
        if err != nil {
            return err
        }
    }

    comparisonStats, err := storage.GetStatsForADate(comparisonDate)
    if err != nil {
		h.logger.Error(err.Error(), "comparison stats, utils: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        if err != sql.ErrNoRows {
            return err
        }
	    h.logger.Info("falling back to default comparison stats", "handler: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        comparisonStats = types.AggStats{
            NumOfLines: currentStats.NumOfLines,
            NumOfCommits: currentStats.NumOfCommits,
            NumOfFiles: currentStats.NumOfFiles,
            NumOfRepos: currentStats.NumOfRepos,
        }
    }
    fmt.Println(currentStats)

    currentStats.DeltaFiles = currentStats.NumOfFiles - comparisonStats.NumOfFiles
    currentStats.DeltaLines = currentStats.NumOfLines - comparisonStats.NumOfLines
    currentStats.DeltaRepos = currentStats.NumOfRepos - comparisonStats.NumOfRepos
    currentStats.DeltaCommits = currentStats.NumOfCommits - comparisonStats.NumOfCommits

	jsonReponse := types.JsonResponse[types.AggStats]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully retrieved the stats",
		Data:          currentStats,
	}
	h.logger.Info("request completed", "handler: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
	utils.WriteJson(w, http.StatusOK, jsonReponse)
	return nil
}

func (h *RepoHandler) UpdateLastCommit(w http.ResponseWriter, r *http.Request) error {
    h.logger.Info("incoming request", "handler: UpdateLastCommit", r.Method, r.URL.Path, r.UserAgent(), r.Body)
    payload, err := utils.ReadJson[types.UpdateLastCommitReq](r)
    if err != nil {
        h.logger.Error(err.Error(), "hander: UpdateLastCommit; Error in parsing request", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        return err
    }

    lastCommitAt, err := strconv.ParseInt(payload.LastCommitAt, 10, 64)
    fmt.Println(lastCommitAt)
    commitTime := time.Unix(lastCommitAt, 0).UTC()
    if err != nil {
        h.logger.Error(err.Error(), "hander: UpdateLastCommit; Invalid date format", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        return err
    }
    fmt.Println(commitTime)

    err = storage.UpdateLastCommit(payload.RepoName, payload.AuthorName, commitTime)
    if err != nil {
        h.logger.Error(err.Error(), "hander: UpdateLastCommit; Failed to update timestamp", r.Method, r.URL.Path, r.UserAgent(), r.Body)
        return err
    }

	jsonResponse := types.JsonResponse[string]{
		RequestStatus: 1,
		StatusCode:    http.StatusOK,
		Msg:           "Successfully updated last commit timestamp",
		Data:          "Commit time updated successfully",
	}
	utils.WriteJson(w, http.StatusOK, jsonResponse)
	h.logger.Info("request completed", "handler: GetStats", r.Method, r.URL.Path, r.UserAgent(), r.Body)
    return nil
}
