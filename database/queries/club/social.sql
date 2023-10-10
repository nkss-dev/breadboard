-- name: CreateClubSocial :exec
INSERT INTO club_social (
    club_name, platform_type, link
)
VALUES (
    @club_name, @platform_type, @link
);

-- name: GetClubSocials :many
SELECT
    platform_type,
    link
FROM
    club_social
WHERE
    club_name = @club_name;

-- name: UpdateClubSocials :exec
UPDATE
    club_social
SET
    link = @link
WHERE
    platform_type = @platform_type
    AND club_name = @club_name;

-- name: DeleteClubSocial :exec
DELETE FROM club_social
WHERE
    club_name = @club_name
    AND platform_type = @platform_type;
