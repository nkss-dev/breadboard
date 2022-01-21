package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "github.com/gorilla/mux"
)

var db *gorm.DB
var err error

type Student struct {
    RollNumber       int    `json:"Roll_Number"`
    Section          string `json:"Section"`
    SubSection       string `json:"SubSection"`
    Name             string `json:"Name"`
    Gender           string `json:"Gender"`
    Mobile           string `json:"Mobile"`
    Birthday         string `json:"Birthday"`
    Email            string `json:"Institute_Email"`
    Batch            int    `json:"Batch"`
    HostelNumber     string `json:"Hostel_Number"`
    RoomNumber       string `json:"Room_Number"`
    DiscordUID       int    `json:"Discord_UID"`
    Verified         bool   `json:"Verified"`
}

func GetUserByRoll(w http.ResponseWriter, r *http.Request) {
    db, err = gorm.Open(sqlite.Open("details.db"), &gorm.Config{})
    if err != nil {
        fmt.Println(err.Error())
    }

    var student Student
    vars := mux.Vars(r)

    db.Table("main").Where("Roll_Number = ?", vars["roll"]).First(&student)
    json.NewEncoder(w).Encode(student)
}
