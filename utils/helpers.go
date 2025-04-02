package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitbook/app/types"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

const MAX_FILE_SIZE = 2 * 1024 * 1024 // 2 MB

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func ReadJson[T any](r *http.Request) (T, error) {
    var reqData T

    const maxBodySize = 1 << 20 // 1MB
    r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&reqData); err != nil {
        if errors.Is(err, io.EOF) {
            return reqData, errors.New("empty request body")
        }
        return reqData, err
    }

    if decoder.More() {
        return reqData, errors.New("request body contains extra data")
    }

    return reqData, nil
}

func RunCommand(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
    fileSize := len(output)
    if fileSize > MAX_FILE_SIZE {
        return "", fmt.Errorf("file is too large to process (%d bytes)", fileSize)
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
