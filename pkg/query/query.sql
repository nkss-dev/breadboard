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
    CAST(ARRAY(SELECT gm.roll_number from group_member gm where g.name = gm.group_name) AS bigint[]) AS members
FROM
    groups g
    JOIN group_discord gd ON g.name = gd.name;
