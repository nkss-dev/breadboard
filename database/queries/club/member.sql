-- name: CreateClubMember :exec
INSERT INTO club_member (
    club_name, roll_number, position, extra_groups, comments
)
VALUES (
    (SELECT c.name FROM club AS c WHERE c.name = @club_name_or_alias OR c.alias = @club_name_or_alias),
    @roll_number,
    @position,
    @extra_groups,
    @comments
);

-- name: ReadClubMembers :many
SELECT
    student.roll_number,
    student.name,
    student.section,
    student.batch,
    student.email,
    club_member.position,
    club_member.extra_groups,
    COALESCE(club_member.comments, '')
FROM
    student
    JOIN club_member ON student.roll_number = club_member.roll_number
WHERE
    club_member.club_name = (SELECT c.name FROM club AS c WHERE c.name = @club_name_or_alias OR c.alias = @club_name_or_alias);

-- name: UpdateClubMember :exec
UPDATE
    club_member
SET
    position = @position,
    extra_groups = @extra_groups,
    comments = @comments
WHERE
    roll_number = @roll_number
    AND club_name = (SELECT c.name FROM club AS c WHERE c.name = @club_name_or_alias OR c.alias = @club_name_or_alias);

-- name: DeleteClubMember :exec
DELETE FROM
    club_member
WHERE
    club_member.club_name = (SELECT c.name FROM club AS c WHERE c.name = @club_name_or_alias OR c.alias = @club_name_or_alias)
    AND club_member.roll_number = @roll_number;
