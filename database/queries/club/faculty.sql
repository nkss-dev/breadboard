-- name: CreateClubFaculty :exec
INSERT INTO club_faculty (
    club_name, emp_id
)
VALUES (
    @club_name, @emp_id
);

-- name: GetClubFaculty :many
SELECT
    f.name, f.mobile
FROM
    faculty AS f
    JOIN club_faculty AS cf ON f.emp_id = cf.emp_id
WHERE
    cf.club_name = @club_name;

-- name: DeleteClubFaculty :exec
DELETE FROM
    club_faculty AS cf
WHERE
    cf.club_name = @club_name
    AND cf.emp_id = @emp_id;
