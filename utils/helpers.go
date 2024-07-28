package utils

import (
	"encoding/json"
	"fmt"
	"gitbook/app/types"
	"log/slog"
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

type Logger interface {
	Info(msg, source, method, path, agent string, data any)
	Error(msg, source, method, path, agent string, data any)
}

type SlogLogger struct{}

func (l *SlogLogger) Info(msg, source, method, path, agent string, data any) {
	slog.Info(
		msg,
		"source", source,
		"method", method,
		"path", path,
		"user_agent", agent,
		"data", data,
	)
}

func (l *SlogLogger) Error(msg, source, method, path, agent string, data any) {
	slog.Error(
		msg,
		"source", source,
		"method", method,
		"path", path,
		"user_agent", agent,
		"data", data,
	)
}
