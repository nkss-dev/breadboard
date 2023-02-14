-- name: GetDiscordLinkStatus :one
SELECT is_verified FROM student WHERE discord_id = $1;

-- name: GetHostels :many
SELECT hostel.*, JSON_AGG(JSON_BUILD_OBJECT('name', warden.name, 'mobile', warden.mobile)) AS "wardens"
FROM hostel
LEFT JOIN warden ON warden.hostel_id = hostel.id
GROUP BY hostel.id;

-- name: GetStudent :one
SELECT
    *
FROM
    student
WHERE roll_number = $1;

-- name: GetStudentByDiscordID :one
SELECT
    *
FROM
    student
WHERE discord_id = $1;
