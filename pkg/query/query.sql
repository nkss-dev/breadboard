-- name: GetStudent :one
SELECT * FROM student
WHERE roll_number = $1;

-- name: GetAllCourses :many
SELECT * FROM course;

-- name: GetBranchCourses :many
SELECT * FROM course
WHERE branch = $1;

-- name: GetCourse :one
SELECT * FROM course
WHERE code = $1 LIMIT 1;

-- name: GetCourses :many
SELECT * FROM course
WHERE branch = $1 AND semester = $2;

-- name: GetSemesterCourses :many
SELECT * FROM course
WHERE semester = $1;
