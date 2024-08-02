package utils

import (
	"encoding/json"
	"fmt"
	"gitbook/app/types"
	"net/http"
	"os/exec"
	"strings"
)

type HTTPError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func RaiseHTTPError(msg string, statusCode int) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func HandlerWrapper(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if apiError, ok := err.(HTTPError); ok {
				jsonResponse := types.JsonResponse[map[string]string]{
					RequestStatus: 0,
					StatusCode:    apiError.StatusCode,
					Msg:           apiError.Msg,
					Data:          map[string]string{},
				}
				WriteJson(w, apiError.StatusCode, jsonResponse)
			}
		}
	}
}

func WriteJson(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func RunCommand(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	outputStr := strings.TrimSuffix(string(output), "\n")
	return outputStr, nil
}

func ExtractRepoPath(path string) string {
    parts := strings.Split(path, "/")
    for i, part := range parts {
        if part == "metadata" && i+2 < len(parts) {
            path := strings.Join(parts[i+2:], "/")
            if path == "" {
                return "."
            }
            return strings.Join(parts[i+2:], "/")
        }
    }
    return "."
}
