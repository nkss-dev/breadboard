-- name: GetHostels :many
SELECT hostel.*, JSON_AGG(JSON_BUILD_OBJECT('name', warden.name, 'mobile', warden.mobile)) AS "wardens"
FROM hostel
LEFT JOIN warden ON warden.hostel_id = hostel.id
GROUP BY hostel.id;

-- name: GetStudent :one
SELECT
    *,
    CAST(ARRAY(SELECT cm.club_name FROM club_member AS cm WHERE cm.roll_number = $1) AS VARCHAR[]) AS clubs
FROM
    student
WHERE roll_number = $1;
