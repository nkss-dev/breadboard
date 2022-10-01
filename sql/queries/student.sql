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
