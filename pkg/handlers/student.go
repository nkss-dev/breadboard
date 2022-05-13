package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Student struct {
	RollNumber   int    `json:"roll_number"`
	Section      string `json:"section"`
	SubSection   string `json:"sub_section"`
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	Mobile       string `json:"mobile"`
	Birthday     string `json:"birthday"`
	Email        string `json:"email"`
	Batch        int    `json:"batch"`
	HostelNumber string `json:"hostel_number"`
	RoomNumber   string `json:"room_number"`
	DiscordUID   int    `json:"discord_uid"`
	Verified     bool   `json:"verified"`
}

func GetStudentByRoll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var student Student
	vars := mux.Vars(r)

	db.Table("student").Where("roll_number = ?", vars["roll"]).First(&student)
	respondJSON(w, 200, student)
}
