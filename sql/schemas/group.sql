CREATE TABLE IF NOT EXISTS club (
    name         VARCHAR(50)   PRIMARY KEY,
    alias        VARCHAR(10)   UNIQUE,
    branch       VARCHAR(3)[]  UNIQUE,
    kind         VARCHAR(17)   NOT NULL CHECK(kind in ('cultural club', 'technical club', 'technical society')),
    description  VARCHAR       NOT NULL
);

CREATE TABLE IF NOT EXISTS club_admin (
    club_name    VARCHAR(50) NOT NULL REFERENCES club(name),
    position     VARCHAR(20) NOT NULL,
    roll_number  CHAR(8)     PRIMARY KEY REFERENCES student(roll_number)
);

CREATE TABLE IF NOT EXISTS club_discord (
    club_name       VARCHAR(50) PRIMARY KEY REFERENCES club(name),
    guild_id        BIGINT      UNIQUE NOT NULL REFERENCES guild(id),
    freshman_role   BIGINT      UNIQUE NOT NULL,
    sophomore_role  BIGINT      UNIQUE NOT NULL,
    junior_role     BIGINT      UNIQUE NOT NULL,
    senior_role     BIGINT      UNIQUE NOT NULL,
    guest_role      BIGINT      UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS club_faculty (
    club_name   VARCHAR(50) REFERENCES club(name),
    emp_id      INT         REFERENCES faculty(emp_id),
    PRIMARY KEY (club_name, emp_id)
);

CREATE TABLE IF NOT EXISTS club_member (
    club_name    VARCHAR(50) REFERENCES club(name),
    roll_number  CHAR(8)     REFERENCES student(roll_number),
    PRIMARY KEY (club_name, roll_number)
);

CREATE TABLE IF NOT EXISTS club_social (
    club_name      VARCHAR(50) REFERENCES club(name),
    platform_type  VARCHAR(15),
    link           VARCHAR     NOT NULL,
    PRIMARY KEY (club_name, platform_type)
);

CREATE OR REPLACE VIEW club_discord_user AS
    SELECT
        s.batch,
        s.discord_id,
        c.name, c.alias,
        cd.guild_id, cs.link,
        cd.freshman_role, cd.sophomore_role, cd.junior_role, cd.senior_role,
        cd.guest_role
    FROM
        club_member AS cm
    JOIN club AS c ON c.name = cm.club_name
    JOIN club_discord AS cd ON cd.club_name = cm.club_name
    JOIN club_social AS cs ON cs.club_name = cm.club_name AND cs.platform_type = 'discord'
    JOIN student AS s ON s.roll_number = cm.roll_number
    WHERE
        s.discord_id IS NOT NULL;
