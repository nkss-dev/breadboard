-- name: GetDiscordLinkStatus :one
SELECT is_verified FROM student WHERE discord_id = $1;

-- name: GetHostels :many
SELECT hostel.*, JSON_AGG(JSON_BUILD_OBJECT('name', warden.name, 'mobile', warden.mobile)) AS "wardens"
FROM hostel
LEFT JOIN warden ON warden.hostel_id = hostel.id
GROUP BY hostel.id;

-- name: GetStudent :one
SELECT
    *,
    (
        SELECT
            COALESCE(JSONB_AGG(JSONB_BUILD_OBJECT(
                'name', cm.club_name,
                'alias', club.alias,
                'position', cm.position,
                'extra_groups', cm.extra_groups
            ) ORDER BY cm.club_name), '[]')::JSONB
        FROM
            club_member AS cm
            JOIN club ON club.name = cm.club_name
        WHERE
            cm.roll_number = student.roll_number
    ) AS clubs
FROM
    student
WHERE student.roll_number = $1;

-- name: GetStudentByDiscordID :one
SELECT
    *
FROM
    student
WHERE discord_id = $1;
