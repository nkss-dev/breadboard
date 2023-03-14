-- name: CreateClubAdmin :exec
WITH new_admin AS (
    UPDATE
        club_details
    SET
        admins = ARRAY_APPEND(admins, @roll_number)
    WHERE
        club_details.club_name = (SELECT club.name FROM club WHERE club.name = @name OR club.alias = @name)
)
UPDATE
    student
SET
    clubs = clubs || JSONB_BUILD_OBJECT(@name::VARCHAR, @position::VARCHAR)
WHERE
    roll_number = @roll_number::CHAR(8);

-- name: CreateClubFaculty :exec
INSERT INTO club_faculty (
    club_name, emp_id
)
VALUES (
    (SELECT c.name from club c WHERE c.name = $1 or c.alias = $1),
    $2
);

-- name: CreateClubMember :exec
INSERT INTO club_member (
    club_name, roll_number
)
VALUES (
    (SELECT name from club WHERE name = $1 or alias = $1),
    $2
);

-- name: CreateClubSocial :exec
INSERT INTO club_social (
    name, platform_type, link
)
VALUES (
    (SELECT c.name from club c WHERE c.name = $1 or c.alias = $1),
    $2,
    $3
);

-- name: DeleteClubAdmin :exec
WITH delete_admin AS (
    UPDATE
        student
    SET
        clubs = clubs - $1
    WHERE
        roll_number = $2
)
UPDATE
    club
SET
    admins = ARRAY_REMOVE(admins, @roll_number::CHAR(8))
WHERE
    club.name = @name
    OR club.alias = @name;

-- name: DeleteClubFaculty :exec
DELETE FROM
    club_faculty AS cf
WHERE
    cf.club_name = (SELECT c.name FROM club c WHERE c.name = $1 OR c.alias = $1)
    AND cf.emp_id = $2;

-- name: DeleteClubMember :exec
DELETE FROM club_member
WHERE
    club_name = (SELECT name FROM club WHERE name = $1 OR alias = $1)
    AND roll_number = $2;

-- name: DeleteClubSocial :exec
DELETE FROM club_social
WHERE
    club_name = (SELECT name FROM club WHERE name = $1 OR alias = $1)
    AND platform_type = $2;

-- name: GetClub :one
SELECT
    club.name,
    COALESCE(club.alias, '') AS alias,
    club.category,
    club.short_description,
    club.email,
    club.is_official,
    COALESCE(JSONB_BUILD_OBJECT(
        'about_us', cd.about_us,
        'why_us', cd.why_us,
        'role_of_sophomore', cd.role_of_soph,
        'role_of_junior', cd.role_of_junior,
        'role_of_senior', cd.role_of_senior
    ), '{}')::JSONB AS description,
    (
        SELECT
            COALESCE(JSONB_AGG(JSONB_BUILD_OBJECT(
                'position', s.clubs -> COALESCE(club.alias, club.name),
                'name', s.name,
                'phone', s.mobile,
                'email', s.email
            ) ORDER BY s.name), '[]')::JSONB
        FROM
            student AS s
        WHERE
            s.roll_number = ANY(cd.admins)
    ) AS admins,
    cd.branch,
    (
        SELECT
            COALESCE(JSONB_AGG(JSONB_BUILD_OBJECT('name', f.name, 'phone', f.mobile) ORDER BY f.name), '[]')::JSONB
        FROM
            faculty AS f
        JOIN club_faculty AS cf ON f.emp_id = cf.emp_id
        WHERE
            cf.club_name = club.name
    ) AS faculties,
    (
        SELECT
            COALESCE(JSONB_AGG(JSONB_BUILD_OBJECT('platform', cs.platform_type, 'link', cs.link) ORDER BY cs.platform_type), '[]')::JSONB
        FROM
            club_social AS cs
        WHERE
            cs.club_name = club.name
    ) AS socials
FROM
    club
JOIN
    club_details AS cd
    ON club.name = cd.club_name
WHERE
    club.name = $1
    OR club.alias = $1;

-- name: GetClubs :many
SELECT
    name,
    COALESCE(alias, name) AS short_name,
    category,
    short_description,
    email,
    is_official
FROM
    club
ORDER BY
    club.name;

-- name: GetClubFaculty :many
SELECT
    f.name, f.mobile
FROM
    faculty AS f
    JOIN club_faculty AS cf ON f.emp_id = cf.emp_id
WHERE
    cf.club_name = $1
    OR $1 = (SELECT alias FROM club WHERE name = cf.club_name);

-- name: GetClubMembers :many
SELECT
    s.*
FROM
    student s
    JOIN club_member member ON s.roll_number = member.roll_number
WHERE
    member.club_name = $1
    OR $1 = (SELECT alias FROM club WHERE name = member.club_name);

-- name: GetClubSocials :many
SELECT
    platform_type,
    link
FROM
    club_social
WHERE
    club_name = $1
    OR $1 = (SELECT alias FROM club WHERE name = club_name);

-- name: UpdateClubSocials :exec
UPDATE
    club_social
SET
    link = $2
WHERE
    platform_type = $1
    AND club_name = $3
    OR $3 = (SELECT alias FROM club WHERE name = club_name);
