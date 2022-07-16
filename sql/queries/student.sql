-- name: GetStudent :one
SELECT * FROM student
WHERE roll_number = $1;

-- name: GetClubMemberships :many
SELECT student.*, group_member.group_name FROM group_member, student WHERE group_member.roll_number = $1 AND group_member.roll_number = student.roll_number;

-- name: GetClubAdmins :many
SELECT student.*, group_admin.group_name FROM group_admin, student WHERE group_admin.roll_number = $1 AND group_admin.roll_number = student.roll_number;
