package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"nkssbackend/query"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

// GetCourse is a handler for retrieving a single course via the `code` argument.
func GetCourse(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		course, err := queries.GetCourse(ctx, vars["code"])
		if err == sql.ErrNoRows {
			RespondError(w, 404, "Course not found in the database")
			return
		}
		RespondJSON(w, 200, course)
	}
}

// GetCourses is a handler for retrieving all the courses matching the given
// query parameters. It outputs all the courses if no parameter is passed.
func GetCourses(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()
		var semester int
		var err error
		var courses []query.Course

		if vars.Get("semester") != "" {
			semester, err = strconv.Atoi(vars.Get("semester"))
			if err != nil {
				RespondError(w, 400, "Semester paramter must be of type int")
				return
			}
			if semester < 1 || semester > 8 {
				RespondError(w, 400, "Semester parameter must be between 1 and 8 (inclusive)")
				return
			}
		}
		branch := vars.Get("branch")
		branches := []string{"", "CE", "CS", "ECE", "EE", "IT", "ME", "PIE"}
		if !slices.Contains(branches, branch) {
			RespondError(w, 400, "Invalid branch parameter!")
			return
		}

		if semester == 0 && branch == "" {
			courses, err = queries.GetAllCourses(ctx)
		} else if semester != 0 && branch == "" {
			courses, err = queries.GetSemesterCourses(ctx, int16(semester))
		} else if semester == 0 && branch != "" {
			courses, err = queries.GetBranchCourses(ctx, branch)
		} else {
			courses, err = queries.GetCourses(ctx, query.GetCoursesParams{Branch: branch, Semester: int16(semester)})
		}
		if err == sql.ErrNoRows || len(courses) == 0 {
			RespondError(w, 404, "Courses not found in the database")
			return
		}
		RespondJSON(w, 200, courses)
	}
}
