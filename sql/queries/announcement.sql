-- name: GetLatestAnnouncementDate :one
SELECT
    MAX(DISTINCT date_of_creation)::DATE
FROM
    academic_announcement;

-- name: GetAcademicAnnouncements :many
SELECT
    date_of_creation,
    title,
    title_link,
    kind
FROM
    academic_announcement
ORDER BY
    date_of_creation DESC;
