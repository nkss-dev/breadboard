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

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", helloWorld).Methods("GET")

    // student.go
    myRouter.HandleFunc("/user/roll/{roll}", GetStudentByRoll).Methods("GET")

    // course.go
    myRouter.HandleFunc("/courses", GetCourses).Methods("GET")
    myRouter.HandleFunc("/course/{code}", GetCourse).Methods("GET")

    // announcement.go
    myRouter.HandleFunc("/announcements", GetAnnouncements).Methods("GET")

    port, port_exists := os.LookupEnv("PORT")
    if !port_exists {
        port = "8081"
    }
    log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func main() {
    fmt.Println("API Online")
    handleRequests()
}
