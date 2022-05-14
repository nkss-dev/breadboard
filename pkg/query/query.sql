-- name: GetStudent :one
SELECT * FROM student
WHERE roll_number = $1;
