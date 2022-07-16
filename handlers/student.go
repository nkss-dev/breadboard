package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"nkssbackend/internal/query"

	"github.com/gorilla/mux"
)

func GetStudentByRoll(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, err := strconv.Atoi(vars["roll"])
		if err != nil {
			RespondError(w, 400, "Roll number parameter must be of type int")
			return
		}

		student, err := queries.GetStudent(ctx, vars["roll"])
		if err == sql.ErrNoRows {
			RespondError(w, 404, "Roll number not found in the database")
			return
		}
		RespondJSON(w, 200, student)
	}
}

func GetStudentClubMemberships(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, err := strconv.Atoi(vars["roll"])
		if err != nil {
			RespondError(w, 400, "Roll number parameter must be of type int")
			return
		}

		student, err := queries.GetClubMemberships(ctx, vars["roll"])
		if err == sql.ErrNoRows || len(student) == 0 {
			RespondError(w, 404, "Student is not in a club")
			return
		}
		RespondJSON(w, 200, student)
	}
}

func IsStudentAdmin(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, err := strconv.Atoi(vars["roll"])
		if err != nil {
			RespondError(w, 400, "Roll number parameter must be of type int")
			return
		}

		student, err := queries.GetClubAdmins(ctx, vars["roll"])
		if err == sql.ErrNoRows || len(student) == 0 {
			RespondError(w, 404, "Student is not an admin")
			return
		}
		RespondJSON(w, 200, student)
	}
}
