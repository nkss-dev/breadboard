-- name: CreateGroupAdmin :exec
INSERT INTO group_admin (
    group_name, position, roll_number
)
VALUES (
    (SELECT name from groups WHERE name = $1 or alias = $1),
    $2,
    $3
);

-- name: CreateGroupFaculty :exec
INSERT INTO group_faculty (
    group_name, emp_id
)
VALUES (
    (SELECT g.name from groups g WHERE g.name = $1 or g.alias = $1),
    $2
);

-- name: CreateGroupMember :exec
INSERT INTO group_member (
    group_name, roll_number
)
VALUES (
    (SELECT name from groups WHERE name = $1 or alias = $1),
    $2
);

-- name: CreateGroupSocial :exec
INSERT INTO group_social (
    name, platform_type, link
)
VALUES (
    (SELECT g.name from groups g WHERE g.name = $1 or g.alias = $1),
    $2,
    $3
);

-- name: DeleteGroupAdmin :exec
DELETE FROM group_admin
WHERE
    group_name = (SELECT name FROM groups WHERE name = $1 OR alias = $1)
    AND roll_number = $2;

-- name: DeleteGroupFaculty :exec
DELETE FROM group_faculty gf
WHERE
    gf.group_name = (SELECT g.name FROM groups g WHERE g.name = $1 OR g.alias = $1)
    AND gf.emp_id = $2;

-- name: DeleteGroupMember :exec
DELETE FROM group_member
WHERE
    group_name = (SELECT name FROM groups WHERE name = $1 OR alias = $1)
    AND roll_number = $2;

-- name: DeleteGroupSocial :exec
DELETE FROM group_social
WHERE
    group_name = (SELECT name FROM groups WHERE name = $1 OR alias = $1)
    AND platform_type = $2;

-- name: GetGroup :one
SELECT
    g.*,
    CAST(ARRAY(SELECT f.name FROM faculty AS f JOIN group_faculty AS gf ON f.emp_id = gf.emp_id WHERE g.name = gf.group_name) AS text[]) AS faculty_names,
    CAST(ARRAY(SELECT f.mobile FROM faculty AS f JOIN group_faculty AS gf ON f.emp_id = gf.emp_id WHERE g.name = gf.group_name) AS text[]) AS faculty_mobiles,
    CAST(ARRAY(SELECT gs.platform_type FROM group_social AS gs WHERE g.name = gs.group_name) AS text[]) AS social_types,
    CAST(ARRAY(SELECT gs.link FROM group_social AS gs WHERE g.name = gs.group_name) AS text[]) AS social_links,
    CAST(ARRAY(SELECT ga.position FROM group_admin AS ga WHERE g.name = ga.group_name) AS text[]) AS admin_positions,
    CAST(ARRAY(SELECT ga.roll_number FROM group_admin AS ga WHERE g.name = ga.group_name) AS bigint[]) AS admin_rolls,
    CAST(ARRAY(SELECT gm.roll_number FROM group_member AS gm WHERE g.name = gm.group_name) AS bigint[]) AS members
FROM
    groups AS g
WHERE
    g.name = $1
    OR g.alias = $1;

-- name: GetGroups :many
SELECT
    g.*,
    CAST(ARRAY(SELECT f.name FROM faculty AS f JOIN group_faculty AS gf ON f.emp_id = gf.emp_id WHERE g.name = gf.group_name) AS text[]) AS faculty_names,
    CAST(ARRAY(SELECT f.mobile FROM faculty AS f JOIN group_faculty AS gf ON f.emp_id = gf.emp_id WHERE g.name = gf.group_name) AS text[]) AS faculty_mobiles,
    CAST(ARRAY(SELECT gs.platform_type FROM group_social AS gs WHERE g.name = gs.group_name) AS text[]) AS social_types,
    CAST(ARRAY(SELECT gs.link FROM group_social AS gs WHERE g.name = gs.group_name) AS text[]) AS social_links,
    CAST(ARRAY(SELECT ga.position FROM group_admin AS ga WHERE g.name = ga.group_name) AS text[]) AS admin_positions,
    CAST(ARRAY(SELECT ga.roll_number FROM group_admin AS ga WHERE g.name = ga.group_name) AS bigint[]) AS admin_rolls,
    CAST(ARRAY(SELECT gm.roll_number FROM group_member AS gm WHERE g.name = gm.group_name) AS bigint[]) AS members
FROM
    groups AS g;

-- name: GetGroupAdmins :many
SELECT
    s.*, admin.position
FROM
    student s
    JOIN group_admin admin ON s.roll_number = admin.roll_number
WHERE
    admin.group_name = $1
    OR $1 = (SELECT alias FROM groups WHERE name = admin.group_name);

-- name: GetGroupFaculty :many
SELECT
    f.name, f.mobile
FROM
    faculty AS f
    JOIN group_faculty AS gf ON f.emp_id = gf.emp_id
WHERE
    gf.group_name = $1
    OR $1 = (SELECT alias FROM groups WHERE name = gf.group_name);

-- name: GetGroupMembers :many
SELECT
    s.*
FROM
    student s
    JOIN group_member member ON s.roll_number = member.roll_number
WHERE
    member.group_name = $1
    OR $1 = (SELECT alias FROM groups WHERE name = member.group_name);

-- name: GetGroupSocials :many
SELECT
    platform_type,
    link
FROM
    group_social
WHERE
    group_name = $1
    OR $1 = (SELECT alias FROM groups WHERE name = group_name);

-- name: UpdateGroupSocials :exec
UPDATE
    group_social
SET
    link = $2
WHERE
    platform_type = $1
    AND group_name = $3
    OR $3 = (SELECT alias FROM groups WHERE name = group_name);
