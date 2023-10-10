-- name: CreateClubMember :exec
INSERT INTO club_member (
    club_name, roll_number, position, extra_groups, comments
)
VALUES (
    @club_name, @roll_number, @position, @extra_groups, @comments
);

-- name: CreateClubMemberBulk :copyfrom
INSERT INTO club_member (
    club_name, roll_number, position, extra_groups, comments
)
VALUES (
    @club_name, @roll_number, @position, @extra_groups, @comments
);

-- name: ReadClubMembers :many
SELECT
    student.roll_number,
    student.section,
    student.name,
    COALESCE(student.mobile, '') AS phone,
    student.email,
    student.batch,
    club_member.position,
    club_member.extra_groups,
    COALESCE(club_member.comments, '')
FROM
    student
    JOIN club_member ON student.roll_number = club_member.roll_number
WHERE
    club_member.club_name = @club_name;

-- name: UpdateClubMember :exec
UPDATE
    club_member
SET
    position = @position,
    extra_groups = @extra_groups,
    comments = @comments
WHERE
    roll_number = @roll_number
    AND club_name = @club_name;

-- name: DeleteClubMember :exec
DELETE FROM
    club_member
WHERE
    club_member.club_name = @club_name
    AND club_member.roll_number = @roll_number;

-- name: DeleteClubMemberBulk :exec
DELETE FROM
    club_member
WHERE
    club_name = @club_name
    AND roll_number = ANY(@roll_numbers::VARCHAR(9)[]);
