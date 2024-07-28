package utils

import (
	"encoding/json"
	"fmt"
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

func RaiseHTTPError(err error, msg string, statusCode int) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func HandlerWrapper(a APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			if apiError, ok := err.(HTTPError); ok {
				WriteJson(w, apiError.StatusCode, apiError)
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
