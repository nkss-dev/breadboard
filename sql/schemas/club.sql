CREATE TABLE IF NOT EXISTS club (
    name               VARCHAR(64)   PRIMARY KEY,
    alias              VARCHAR(16)   UNIQUE,
    category           VARCHAR(32)   NOT NULL,
    email              VARCHAR(64)   NOT NULL,
    short_description  VARCHAR(256)  NOT NULL,
    is_official        BOOLEAN       DEFAULT false NOT NULL,
    CONSTRAINT ck_category CHECK (
        category IN (
            'Committee',
            'Crew',
            'Cultural Club',
            'Magazine',
            'Technical Club',
            'Technical Society'
        )
    )
);

CREATE TABLE IF NOT EXISTS club_details (
    club_name       VARCHAR(64)  PRIMARY KEY REFERENCES club(name) ON UPDATE CASCADE,
    about_us        VARCHAR      NOT NULL,
    why_us          VARCHAR      NOT NULL,
    role_of_soph    VARCHAR      NOT NULL,
    role_of_junior  VARCHAR      NOT NULL,
    role_of_senior  VARCHAR      NOT NULL,
    admins          CHAR(8)[]    NOT NULL,
    branch          CHAR(2)[]    NOT NULL
);

CREATE TABLE IF NOT EXISTS club_discord (
    club_name       VARCHAR(64) PRIMARY KEY REFERENCES club(name) ON UPDATE CASCADE,
    guild_id        BIGINT      UNIQUE NOT NULL REFERENCES guild(id),
    guest_role      BIGINT      UNIQUE,
    member_role     BIGINT      UNIQUE
);

CREATE TABLE IF NOT EXISTS club_faculty (
    club_name   VARCHAR(64) REFERENCES club(name) ON UPDATE CASCADE,
    emp_id      INT         REFERENCES faculty(emp_id),
    PRIMARY KEY (club_name, emp_id)
);

CREATE TABLE IF NOT EXISTS club_member (
    club_name    VARCHAR(64) REFERENCES club(name) ON UPDATE CASCADE,
    roll_number  CHAR(8)     REFERENCES student(roll_number),
    PRIMARY KEY (club_name, roll_number)
);

CREATE TABLE IF NOT EXISTS club_social (
    club_name      VARCHAR(64) REFERENCES club(name) ON UPDATE CASCADE,
    platform_type  VARCHAR(15),
    link           VARCHAR     NOT NULL,
    PRIMARY KEY (club_name, platform_type)
);

CREATE OR REPLACE VIEW club_discord_user AS
    SELECT
        s.batch,
        s.discord_id,
        c.name, c.alias,
        cd.guild_id, (SELECT link FROM club_social WHERE platform_type = 'discord') AS guild_invite,
        cd.guest_role,
        cd.member_role
    FROM
        club_member AS cm
    JOIN club AS c ON c.name = cm.club_name
    JOIN club_discord AS cd ON cd.club_name = cm.club_name
    JOIN student AS s ON s.roll_number = cm.roll_number
    WHERE
        s.discord_id IS NOT NULL;
