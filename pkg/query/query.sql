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

-- name: GetAllGroups :many
SELECT * FROM groups NATURAL JOIN group_discord;

-- name: GetAllFaculty :many
SELECT * FROM group_faculty;

-- name: GetAllGroupSocials :many
SELECT * FROM group_social;

-- name: GetAllGroupAdmins :many
SELECT * FROM group_admin;

-- name: GetAllGroupMembers :many
SELECT * FROM group_member;
