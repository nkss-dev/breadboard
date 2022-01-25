package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", helloWorld).Methods("GET")
    myRouter.HandleFunc("/user/roll/{roll}", GetStudentByRoll).Methods("GET")
    myRouter.HandleFunc("/courses", GetCourses).Methods("GET")
    myRouter.HandleFunc("/course/{code}", GetCourse).Methods("GET")
    log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
    fmt.Println("API Online")
    handleRequests()
}
