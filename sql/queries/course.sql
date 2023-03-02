-- name: CreateCourse :exec
INSERT INTO course (
    code,
    title,
    prereq,
    kind,
    objectives,
    content,
    book_names,
    outcomes
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (code)
DO
    UPDATE SET
        title = $2,
        prereq = $3,
        kind = $4,
        objectives = $5,
        content = $6,
        book_names = $7,
        outcomes = $8;

-- name: CreateSpecifics :exec
INSERT INTO branch_specifics (
    code,
    branch,
    semester,
    credits
)
VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (code, branch)
DO
    UPDATE SET
        semester = $3,
        credits = $4;

-- name: GetCourse :one
SELECT
    c.*, (
        SELECT
            COALESCE(JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            )), '[]')::JSON
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code = $1;

-- name: GetCourses :many
SELECT
    c.*, (
        SELECT
            COALESCE(JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            )), '[]')::JSON
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c;

-- name: GetCoursesByBranch :many
SELECT
    c.*, (
        SELECT
            JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            ))
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code IN (SELECT DISTINCT bs.code FROM branch_specifics AS bs WHERE bs.branch = $1);

-- name: GetCoursesByBranchAndSemester :many
SELECT
    c.*, (
        SELECT
            JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            ))
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code IN (SELECT DISTINCT bs.code FROM branch_specifics AS bs WHERE bs.branch = $1 AND bs.semester = $2);

-- name: GetCoursesBySemester :many
SELECT
    c.*, (
        SELECT
            JSON_AGG(JSON_BUILD_OBJECT(
                'branch', bs.branch,
                'semester', bs.semester,
                'credits', bs.credits
            ))
        FROM
            branch_specifics AS bs
        WHERE
            bs.code = c.code
    ) AS specifics
FROM
    course AS c
WHERE
    c.code IN (SELECT DISTINCT bs.code FROM branch_specifics AS bs WHERE bs.semester = $1);
