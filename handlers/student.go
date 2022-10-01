package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"nkssbackend/internal/query"

	"github.com/gorilla/mux"
)

// GetHostels retrieves all the hostels and their meta data from the database.
func GetHostels(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		hostels, _ := queries.GetHostels(ctx)
		RespondJSON(w, 200, hostels)
	}
}

// GetStudent retrieves a single student's details based on either their roll number or Discord ID.
func GetStudent(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, err := strconv.Atoi(vars["roll"])
		if err != nil {
			RespondError(w, 400, "Roll number parameter must be an integer")
			return
		}

		student, err := queries.GetStudent(ctx, vars["roll"])
		if err == sql.ErrNoRows {
			RespondError(w, 404, "Student not found in the database")
			return
		}
		RespondJSON(w, 200, student)
	}
}
