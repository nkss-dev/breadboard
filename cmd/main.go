package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    myRouter := mux.NewRouter().StrictSlash(true)

    fmt.Println("API Online")

    myRouter.HandleFunc("/", helloWorld).Methods("GET")

    port, port_exists := os.LookupEnv("PORT")
    if !port_exists {
        port = "8081"
    }
    log.Fatal(http.ListenAndServe(":"+port, myRouter))
}
