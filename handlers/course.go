package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	query "breadboard/.sqlc-auto-gen"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

type BranchSpecifics struct {
	Branch   string
	Semester int16
	Credits  []int16
}

type Course struct {
	Code       string
	Title      string
	Prereq     []string
	Kind       string
	Objectives []string
	Content    []string
	BookNames  []string `json:"book_names"`
	Outcomes   []string
	Specifics  []BranchSpecifics
}

func CreateCourse(conn *pgx.Conn) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		var coursePaths []string
		json.NewDecoder(r.Body).Decode(&coursePaths)

		courseMarkdowns := fetchCourse(coursePaths)

		for _, courseMarkdown := range courseMarkdowns {
			course, err := parseMarkdown(courseMarkdown)
			if err != nil {
				RespondError(w, 400, "Failed to parse the markdown file. Error: "+err.Error())
				return
			}

			var queryParam = query.CreateCourseParams{
				Code:       course.Code,
				Title:      course.Title,
				Prereq:     course.Prereq,
				Kind:       course.Kind,
				Objectives: course.Objectives,
				Content:    course.Content,
				BookNames:  course.BookNames,
				Outcomes:   course.Outcomes,
			}
			err = queries.CreateCourse(ctx, queryParam)
			if err != nil {
				log.Println(err)
				RespondError(w, 500, "Something went wrong while inserting the course to our database")
				return
			}

			for _, specific := range course.Specifics {
				var queryParam = query.CreateSpecificsParams{
					Code:     course.Code,
					Branch:   specific.Branch,
					Semester: specific.Semester,
					Credits:  specific.Credits,
				}
				err := queries.CreateSpecifics(ctx, queryParam)
				if err != nil {
					log.Println(err)
					RespondError(w, 500, "Something went wrong while inserting the course specifics to our database")
					return
				}
			}
		}
		RespondJSON(w, 201, "")
	}
}

// GetCourse is a handler for retrieving a single course via the `code` argument.
func GetCourse(conn *pgx.Conn) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		course, err := queries.GetCourse(ctx, vars["code"])
		if err == pgx.ErrNoRows {
			RespondError(w, 404, "Course not found in the database")
			return
		}
		RespondJSON(w, 200, course)
	}
}

// GetCourses is a handler for retrieving all the courses matching the given
// query parameters. It outputs all the courses if no parameter is passed.
func GetCourses(conn *pgx.Conn) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()
		var semester int
		var err error

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

		// !TODO: Make code idempotent
		if semester == 0 && branch == "" {
			courses, err := queries.GetCourses(ctx)
			if err != nil && err != pgx.ErrNoRows {
				RespondError(w, 500, err.Error())
				return
			} else if err == pgx.ErrNoRows || len(courses) == 0 {
				RespondError(w, 404, "Courses not found in the database")
				return
			}
			RespondJSON(w, 200, courses)
		} else if semester != 0 && branch == "" {
			courses, err := queries.GetCoursesBySemester(ctx, int16(semester))
			if err == pgx.ErrNoRows || len(courses) == 0 {
				RespondError(w, 404, "Courses not found in the database")
				return
			}
			RespondJSON(w, 200, courses)
		} else if semester == 0 && branch != "" {
			courses, err := queries.GetCoursesByBranch(ctx, branch)
			if err == pgx.ErrNoRows || len(courses) == 0 {
				RespondError(w, 404, "Courses not found in the database")
				return
			}
			RespondJSON(w, 200, courses)
		} else {
			courses, err := queries.GetCoursesByBranchAndSemester(ctx, query.GetCoursesByBranchAndSemesterParams{Branch: branch, Semester: int16(semester)})
			if err == pgx.ErrNoRows || len(courses) == 0 {
				RespondError(w, 404, "Courses not found in the database")
				return
			}
			RespondJSON(w, 200, courses)
		}
	}
}

func fetchCourse(coursePaths []string) (courses []string) {
	ctx := context.Background()
	client := github.NewClient(nil)
	owner := "NIT-KKR-Student-Support-System"
	repo := "atlas"

	for _, coursePath := range coursePaths {
		fileContents, _, _, err := client.Repositories.GetContents(ctx, owner, repo, coursePath, nil)
		if err != nil {
			log.Println("Error getting file contents:", err)
			continue
		}
		course, err := fileContents.GetContent()
		if err != nil {
			log.Println("Error getting file contents:", err)
			continue
		}
		courses = append(courses, course)
	}
	return courses
}

func parseMarkdown(md string) (Course, error) {
	var course Course

	parts := strings.SplitN(md, "---\n", 3)
	if len(parts) < 3 {
		return course, fmt.Errorf("invalid markdown format")
	}
	err := yaml.Unmarshal([]byte(parts[1]), &course)
	if err != nil {
		return course, err
	}

	fields := strings.SplitN(parts[2], "\n# ", 5)

	for _, objective := range strings.Split(fields[1], "\n- ")[1:] {
		course.Objectives = append(course.Objectives, strings.Trim(objective, "\n"))
	}
	for _, unit := range strings.Split(fields[2], "## Unit ")[1:] {
		course.Content = append(course.Content, strings.Trim(unit[1:], "\n"))
	}
	for _, bookName := range strings.Split(fields[3], "\n- ")[1:] {
		course.BookNames = append(course.BookNames, strings.Trim(bookName, "\n"))
	}
	for _, outcome := range strings.Split(fields[4], "\n- ")[1:] {
		course.Outcomes = append(course.Outcomes, strings.Trim(outcome, "\n"))
	}

	return course, nil
}
