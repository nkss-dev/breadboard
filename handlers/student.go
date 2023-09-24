package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	query "breadboard/.sqlc-auto-gen"

	"github.com/gorilla/mux"
)

func GetDiscordLinkStatus(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			RespondError(w, 400, "ID parameter must be of type int")
			return
		}

		student, err := queries.GetDiscordLinkStatus(ctx, sql.NullInt64{Int64: int64(idInt), Valid: true})
		if err == sql.ErrNoRows {
			RespondJSON(w, 404, false)
		} else {
			RespondJSON(w, 200, student)
		}
	}
}

// GetHostels retrieves all the hostels and their meta data from the database.
func GetHostels(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		hostels, _ := queries.GetHostels(ctx)
		RespondJSON(w, 200, hostels)
	}
}

// GetStudent retrieves a single student's details based on their roll number, email, or Discord ID.
func GetStudent(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)

	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		discordId, _ := strconv.Atoi(id)

		student, err := queries.GetStudent(ctx, query.GetStudentParams{
			RollNumber: id,
			Email:      id,
			DiscordID:  int64(discordId),
		})
		if err == sql.ErrNoRows {
			RespondError(w, 404, "Student not found in the database")
			return
		}
		RespondJSON(w, 200, student)
	}
}
