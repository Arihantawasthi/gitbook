package main

import (
	"log"
	"net/http"
    "encoding/json"
)


func main() {
    router := http.NewServeMux()
    router.HandleFunc("GET api/v1/", func(w http.ResponseWriter, r *http.Request) {
        items := map[string]int{"ROUTE": 1}
        json.NewEncoder(w).Encode(items)
    })

    server := http.Server{
        Addr: ":8000",
        Handler: router,
    }

    log.Printf("Starting server at port %v\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}
