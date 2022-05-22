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
		respondError(w, 400, "Roll number paramter must be of type int")
		return
	}

	ctx := context.Background()
	queries := query.New(db)
	student, err := queries.GetStudent(ctx, int32(roll))
	if err == sql.ErrNoRows {
		respondError(w, 404, "Roll number not found in the database")
		return
	}
	respondJSON(w, 200, student)
}
