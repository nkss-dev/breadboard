package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"NKSS-backend/pkg/query"

	"github.com/gorilla/mux"
)

func GetStudentByRoll(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roll, err := strconv.Atoi(vars["roll"])
	if err != nil {
		RespondError(w, 400, "Roll number parameter must be of type int")
		return
	}

	ctx := context.Background()
	queries := query.New(db)
	student, err := queries.GetStudent(ctx, int32(roll))
	if err == sql.ErrNoRows {
		RespondError(w, 404, "Roll number not found in the database")
		return
	}
	RespondJSON(w, 200, student)
}

func GetStudentClubMemberships(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roll, err := strconv.Atoi(vars["roll"])
	if err != nil {
		RespondError(w, 400, "Roll number parameter must be of type int")
		return
	}

	ctx := context.Background()
	queries := query.New(db)
	student, err := queries.GetClubMemberships(ctx, int32(roll))
	if (err == sql.ErrNoRows || len(student) == 0) {
		RespondError(w, 404, "Student is not in a club")
		return
	}
	RespondJSON(w, 200, student)
}

func IsStudentAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roll, err := strconv.Atoi(vars["roll"])
	if err != nil {
		RespondError(w, 400, "Roll number parameter must be of type int")
		return
	}

	ctx := context.Background()
	queries := query.New(db)
	student, err := queries.GetClubAdmins(ctx, int32(roll))
	if (err == sql.ErrNoRows || len(student) == 0) {
		RespondError(w, 404, "Student is not an admin")
		return
	}
	RespondJSON(w, 200, student)
}
