package app

import (
	"gitbook/app/handler"
    "log/slog"
	"net/http"
)


func HandlerFuncWraper(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming Request", "X-Request-ID", r.Header.Get("X-Request-ID"), "method", r.Method, "path", r.URL.Path, "body", r.Body)
		w.Header().Add("Content-Type", "application/json")
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}


func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /api/v1/repos", HandlerFuncWraper(handler.GetAllRepos))
}
