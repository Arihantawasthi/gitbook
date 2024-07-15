package main

import (
	"log"
	"net/http"
    "gitbook/app"
)


func muxWrap(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
        w.Header().Add("Content-Type", "application/json")
    })
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
