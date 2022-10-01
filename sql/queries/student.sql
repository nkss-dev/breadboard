-- name: GetDiscordLinkStatus :one
SELECT is_verified FROM student WHERE discord_id = $1;

-- name: GetHostels :many
SELECT hostel.*, JSON_AGG(JSON_BUILD_OBJECT('name', warden.name, 'mobile', warden.mobile)) AS "wardens"
FROM hostel
LEFT JOIN warden ON warden.hostel_id = hostel.id
GROUP BY hostel.id;

-- name: GetStudent :one
SELECT
    *, (
        SELECT
            JSON_AGG(JSON_BUILD_OBJECT(
                'name', cm.club_name, 'position',
                COALESCE((SELECT position FROM club_admin WHERE roll_number = cm.roll_number), 'Member')
            ))
        FROM
            club_member AS cm
        WHERE
            cm.roll_number = $1
    ) AS clubs
FROM
    student
WHERE roll_number = $1;

-- name: GetStudentByDiscordID :one
SELECT
    *,
    CAST(ARRAY(SELECT club.alias FROM club JOIN club_member AS cm ON cm.club_name = club.name WHERE cm.roll_number = s.roll_number) AS VARCHAR[]) AS clubs
FROM
    student AS s
WHERE discord_id = $1;
