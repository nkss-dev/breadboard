package main

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var db *gorm.DB
var err error

type Student struct {
    RollNumber   int    `json:"Roll_Number"`
    Section      string `json:"Section"`
    SubSection   string `json:"SubSection"`
    Name         string `json:"Name"`
    Gender       string `json:"Gender"`
    Mobile       string `json:"Mobile"`
    Birthday     string `json:"Birthday"`
    Email        string `json:"Institute_Email"`
    Batch        int    `json:"Batch"`
    HostelNumber string `json:"Hostel_Number"`
    RoomNumber   string `json:"Room_Number"`
    DiscordUID   int    `json:"Discord_UID"`
    Verified     bool   `json:"Verified"`
}

func GetStudentByRoll(w http.ResponseWriter, r *http.Request) {
    db, err = gorm.Open(sqlite.Open("details.db"), &gorm.Config{})
    if err != nil {
        fmt.Println(err.Error())
    }

    var student Student
    vars := mux.Vars(r)

    db.Table("main").Where("Roll_Number = ?", vars["roll"]).First(&student)
    json.NewEncoder(w).Encode(student)
}
