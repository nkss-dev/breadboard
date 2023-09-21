-- name: CreateClubMember :exec
INSERT INTO club_member (
    club_name, roll_number, position, extra_groups
)
VALUES (
    (SELECT c.name FROM club AS c WHERE c.name = @club_name_or_alias OR c.alias = @club_name_or_alias),
    @roll_number,
    @position,
    @extra_groups
);

-- name: DeleteClubMember :exec
DELETE FROM
    club_member
WHERE
    club_member.club_name = (SELECT c.name FROM club AS c WHERE c.name = @club_name_or_alias OR c.alias = @club_name_or_alias)
    AND club_member.roll_number = @roll_number;
