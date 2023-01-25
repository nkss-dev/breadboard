-- name: InsertAcademicAnnouncement :exec
INSERT INTO academic_announcement (
    date_of_creation, title, title_link, kind
)
VALUES (
    $1,
    $2,
    $3,
    'academic'
) ON CONFLICT (date_of_creation, title)
DO NOTHING;

-- name: GetAcademicAnnouncements :exec
SELECT date_of_creation, title, title_link, kind
FROM academic_announcement
ORDER BY date_of_creation DESC;
