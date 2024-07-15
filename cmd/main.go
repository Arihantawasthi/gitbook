package main

import (
	"gitbook/app"
	"gitbook/app/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
)


func muxWrap(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        slog.Info("Incoming Request", "X-Request-ID", r.Header.Get("X-Request-ID"), "method", r.Method, "path", r.URL.Path, "body", r.Body)
        next.ServeHTTP(w, r)
        w.Header().Add("Content-Type", "application/json")
    })
}

func init() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)
    storage.ConnectPGStorage()
}

func main() {
    router := http.NewServeMux()
    app.RegisterRoutes(router)

    server := http.Server{
        Addr: ":8000",
        Handler: muxWrap(router),
    }

    log.Printf("Starting server at port %v\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}
