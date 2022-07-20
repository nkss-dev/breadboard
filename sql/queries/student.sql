-- name: GetStudent :one
SELECT
    *,
    CAST(ARRAY(SELECT cm.club_name FROM club_member AS cm WHERE cm.roll_number = $1) AS VARCHAR[]) AS clubs
FROM
    student
WHERE roll_number = $1;
