package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"nkssbackend/internal/query"

	"github.com/gorilla/mux"
)

func GetHostels(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		hostels, _ := queries.GetHostels(ctx)
		RespondJSON(w, 200, hostels)
	}
}

func GetStudent(db *sql.DB) http.HandlerFunc {
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
