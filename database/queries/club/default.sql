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
                'position', club_member.position,
                'roll', s.roll_number,
                'name', s.name,
                'phone', s.mobile,
                'email', s.email
            ) ORDER BY s.name), '[]')::JSONB
        FROM
            student AS s
            JOIN club_member ON s.roll_number = club_member.roll_number AND club.name = club_member.club_name
        WHERE
            s.roll_number = ANY(SELECT roll_number FROM club_member WHERE club_name = club.name AND position != 'Member')
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
    club.name = @club_name;

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
