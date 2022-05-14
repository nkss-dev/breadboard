package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"NKSS-backend/pkg/query"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

// func ParseCourseMD() map[string]interface{} {
// 	markdown := goldmark.New(
// 		goldmark.WithExtensions(
// 			meta.Meta,
// 		),
// 	)
// 	file, err := ioutil.ReadFile("courses/HSIR11.md")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	var buf bytes.Buffer
// 	if err := goldmark.Convert(file, &buf); err != nil {
// 		panic(err)
// 	}
// 	// fmt.Println(buf.String())
// 	context := parser.NewContext()
// 	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
// 		panic(err)
// 	}
// 	metaData := meta.Get(context)
// 	// title := metaData["for"]
// 	// fmt.Print(metaData)
// 	return metaData
// }

func GetCourse(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ctx := context.Background()
	queries := query.New(db)
	course, err := queries.GetCourse(ctx, vars["code"])
	if err == sql.ErrNoRows {
		respondError(w, 404, "Course not found in the database")
		return
	}
	fmt.Println(course)
	respondJSON(w, 200, course)
}

func GetCourses(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var semester int
	var err error
	var courses []query.Course

	if vars.Get("semester") != "" {
		semester, err = strconv.Atoi(vars.Get("semester"))
		if err != nil {
			respondError(w, 400, "Semester paramter must be of type int")
			return
		}
		if semester < 1 || semester > 8 {
			respondError(w, 400, "Semester parameter must be between 1 and 8 (inclusive)")
			return
		}
	}
	branch := vars.Get("branch")
	branches := []string{"", "CE", "CS", "ECE", "EE", "IT", "ME", "PIE"}
	if !slices.Contains(branches, branch) {
		respondError(w, 400, "Invalid branch parameter!")
		return
	}

	ctx := context.Background()
	queries := query.New(db)

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
		respondError(w, 404, "Courses not found in the database")
		return
	}
	respondJSON(w, 200, courses)
}
