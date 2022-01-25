package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Course struct {
    Code          string
    Title         string
    Branch        string
    Semester      int
    Credits       int
    Prerequisites string
    Type          string
    Objectives    string
    Content       string
    Books         string
    Outcomes      string
}

func GetCourse(w http.ResponseWriter, r *http.Request) {
    db, err := gorm.Open(sqlite.Open("db/courses.db"), &gorm.Config{})
    if err != nil {
        fmt.Println(err.Error())
    }
    w.Header().Set("Content-Type", "application/json")

    var course Course
    vars := mux.Vars(r)

    err = db.Where("code = ?", vars["code"]).First(&course).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        resp := make(map[string]string)
        resp["message"] = "No course named " + vars["code"] + " was found"
        jsonResp, _ := json.Marshal(resp)
        w.WriteHeader(http.StatusBadRequest)
        w.Write(jsonResp)
        return
    }

    jsonResp, _ := json.Marshal(course)
    w.Write(jsonResp)
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
    db, err := gorm.Open(sqlite.Open("db/courses.db"), &gorm.Config{})
    if err != nil {
        fmt.Println(err.Error())
    }
    w.Header().Set("Content-Type", "application/json")

    semStr := r.URL.Query().Get("semester")
    semInt, err := strconv.Atoi(semStr)
    if err != nil {
        resp := make(map[string]string)
        resp["message"] = "Semester value must be an integer"
        jsonResp, _ := json.Marshal(resp)
        w.WriteHeader(http.StatusBadRequest)
        w.Write(jsonResp)
        return
    }
    branch := r.URL.Query().Get("branch")

    var courses []Course
    db.Where(&Course{Semester: semInt, Branch: branch}).Find(&courses)

    if len(courses) == 0 {
        resp := make(map[string]string)
        resp["message"] = "No course matching the given constraints was found"
        jsonResp, _ := json.Marshal(resp)
        w.WriteHeader(http.StatusNotFound)
        w.Write(jsonResp)
        return
    }

    jsonResp, _ := json.Marshal(courses)
    w.Write(jsonResp)
}
