-- name: GetStudent :one
SELECT * FROM student
WHERE roll_number = $1;

-- name: GetAllCourses :many
SELECT * FROM course;

-- name: GetBranchCourses :many
SELECT * FROM course
WHERE branch = $1;

-- name: GetCourse :one
SELECT * FROM course
WHERE code = $1 LIMIT 1;

-- name: GetCourses :many
SELECT * FROM course
WHERE branch = $1 AND semester = $2;

-- name: GetSemesterCourses :many
SELECT * FROM course
WHERE semester = $1;

-- name: GetGroup :one
SELECT
    g.*,
    gd.id AS server_id,
    gd.invite AS server_invite,
    gd.fresher_role,
    gd.sophomore_role,
    gd.junior_role,
    gd.senior_role,
    gd.guest_role,
    CAST(ARRAY(SELECT gf.name FROM group_faculty gf WHERE g.name = gf.group_name) AS text[]) AS faculty_names,
    CAST(ARRAY(SELECT gf.mobile FROM group_faculty gf WHERE g.name = gf.group_name) AS bigint[]) AS faculty_mobiles,
    CAST(ARRAY(SELECT gs.type FROM group_social gs WHERE g.name = gs.name) AS text[]) AS social_types,
    CAST(ARRAY(SELECT gs.link FROM group_social gs WHERE g.name = gs.name) AS text[]) AS social_links,
    CAST(ARRAY(SELECT ga.position FROM group_admin ga WHERE g.name = ga.group_name) AS text[]) AS admin_positions,
    CAST(ARRAY(SELECT ga.roll_number FROM group_admin ga WHERE g.name = ga.group_name) AS bigint[]) AS admin_rolls,
    CAST(ARRAY(SELECT gm.roll_number FROM group_member gm WHERE g.name = gm.group_name) AS bigint[]) AS members
FROM
    groups g
    JOIN group_discord gd ON g.name = gd.name
WHERE
    g.name = $1
    OR g.alias = $1;

-- name: GetAllGroups :many
SELECT
    g.*,
    gd.id AS server_id,
    gd.invite AS server_invite,
    gd.fresher_role,
    gd.sophomore_role,
    gd.junior_role,
    gd.senior_role,
    gd.guest_role,
    CAST(ARRAY(SELECT gf.name FROM group_faculty gf WHERE g.name = gf.group_name) AS text[]) AS faculty_names,
    CAST(ARRAY(SELECT gf.mobile FROM group_faculty gf WHERE g.name = gf.group_name) AS bigint[]) AS faculty_mobiles,
    CAST(ARRAY(SELECT gs.type FROM group_social gs WHERE g.name = gs.name) AS text[]) AS social_types,
    CAST(ARRAY(SELECT gs.link FROM group_social gs WHERE g.name = gs.name) AS text[]) AS social_links,
    CAST(ARRAY(SELECT ga.position FROM group_admin ga WHERE g.name = ga.group_name) AS text[]) AS admin_positions,
    CAST(ARRAY(SELECT ga.roll_number FROM group_admin ga WHERE g.name = ga.group_name) AS bigint[]) AS admin_rolls,
    CAST(ARRAY(SELECT gm.roll_number FROM group_member gm WHERE g.name = gm.group_name) AS bigint[]) AS members
FROM
    groups g
    JOIN group_discord gd ON g.name = gd.name;

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
    name, mobile
FROM
    group_faculty gf
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
    type,
    link
FROM
    group_social gs
WHERE
    gs.name = $1
    OR $1 = (SELECT alias FROM groups WHERE name = gs.name);

-- name: GetClubMemberships :many
SELECT student.*, group_member.group_name FROM group_member, student WHERE group_member.roll_number = $1 AND group_member.roll_number = student.roll_number;

-- name: GetClubAdmins :many
SELECT student.*, group_admin.group_name FROM group_admin, student WHERE group_admin.roll_number = $1 AND group_admin.roll_number = student.roll_number;
