-- name: CreateClubAdmin :exec
INSERT INTO club_admin (
    club_name, position, roll_number
)
VALUES (
    (SELECT name from club WHERE name = $1 or alias = $1),
    $2,
    $3
);

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
DELETE FROM club_admin
WHERE
    club_name = (SELECT name FROM club WHERE name = $1 OR alias = $1)
    AND roll_number = $2;

-- name: DeleteClubFaculty :exec
DELETE FROM club_faculty cf
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
    c.*,
    CAST(ARRAY(SELECT f.name FROM faculty AS f JOIN club_faculty AS cf ON f.emp_id = cf.emp_id WHERE c.name = cf.club_name) AS text[]) AS faculty_names,
    CAST(ARRAY(SELECT f.mobile FROM faculty AS f JOIN club_faculty AS cf ON f.emp_id = cf.emp_id WHERE c.name = cf.club_name) AS text[]) AS faculty_mobiles,
    CAST(ARRAY(SELECT cs.platform_type FROM club_social AS cs WHERE c.name = cs.club_name) AS text[]) AS social_types,
    CAST(ARRAY(SELECT cs.link FROM club_social AS cs WHERE c.name = cs.club_name) AS text[]) AS social_links,
    CAST(ARRAY(SELECT ca.position FROM club_admin AS ca WHERE c.name = ca.club_name) AS text[]) AS admin_positions,
    CAST(ARRAY(SELECT ca.roll_number FROM club_admin AS ca WHERE c.name = ca.club_name) AS bigint[]) AS admin_rolls,
    CAST(ARRAY(SELECT cm.roll_number FROM club_member AS cm WHERE c.name = cm.club_name) AS bigint[]) AS members
FROM
    club AS c
WHERE
    c.name = $1
    OR c.alias = $1;

-- name: GetClubs :many
SELECT
    c.*,
    CAST(ARRAY(SELECT f.name FROM faculty AS f JOIN club_faculty AS cf ON f.emp_id = cf.emp_id WHERE c.name = cf.club_name) AS text[]) AS faculty_names,
    CAST(ARRAY(SELECT f.mobile FROM faculty AS f JOIN club_faculty AS cf ON f.emp_id = cf.emp_id WHERE c.name = cf.club_name) AS text[]) AS faculty_mobiles,
    CAST(ARRAY(SELECT cs.platform_type FROM club_social AS cs WHERE c.name = cs.club_name) AS text[]) AS social_types,
    CAST(ARRAY(SELECT cs.link FROM club_social AS cs WHERE c.name = cs.club_name) AS text[]) AS social_links,
    CAST(ARRAY(SELECT ca.position FROM club_admin AS ca WHERE c.name = ca.club_name) AS text[]) AS admin_positions,
    CAST(ARRAY(SELECT ca.roll_number FROM club_admin AS ca WHERE c.name = ca.club_name) AS bigint[]) AS admin_rolls,
    CAST(ARRAY(SELECT cm.roll_number FROM club_member AS cm WHERE c.name = cm.club_name) AS bigint[]) AS members
FROM
    club AS c;

-- name: GetClubAdmins :many
SELECT
    s.*, admin.position
FROM
    student s
    JOIN club_admin admin ON s.roll_number = admin.roll_number
WHERE
    admin.club_name = $1
    OR $1 = (SELECT alias FROM club WHERE name = admin.club_name);

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
