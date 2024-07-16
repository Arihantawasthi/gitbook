package main

import (
	"gitbook/app"
	"gitbook/app/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
)


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
        Handler: router,
    }

    log.Printf("Starting server at port %v\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}
