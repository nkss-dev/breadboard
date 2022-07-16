CREATE TABLE IF NOT EXISTS groups (
    name         VARCHAR(50)   PRIMARY KEY,
    alias        VARCHAR(10)   UNIQUE,
    branch       VARCHAR(3)[]  UNIQUE,
    kind         VARCHAR(17)   NOT NULL CHECK(kind in ('cultural club', 'technical club', 'technical society')),
    description  VARCHAR       NOT NULL
);

CREATE TABLE IF NOT EXISTS group_admin (
    group_name   VARCHAR(50) NOT NULL REFERENCES groups(name),
    position     VARCHAR(20) NOT NULL,
    roll_number  CHAR(8)     PRIMARY KEY REFERENCES student(roll_number)
);

CREATE TABLE IF NOT EXISTS group_discord (
    group_name      VARCHAR(50) PRIMARY KEY REFERENCES groups(name),
    server_id       BIGINT      UNIQUE NOT NULL REFERENCES guild(id),
    freshman_role   BIGINT      UNIQUE NOT NULL,
    sophomore_role  BIGINT      UNIQUE NOT NULL,
    junior_role     BIGINT      UNIQUE NOT NULL,
    senior_role     BIGINT      UNIQUE NOT NULL,
    guest_role      BIGINT      UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS group_faculty (
    group_name  VARCHAR(50) REFERENCES groups(name),
    emp_id      INT         REFERENCES faculty(emp_id),
    PRIMARY KEY (group_name, emp_id)
);

CREATE TABLE IF NOT EXISTS group_member (
    group_name   VARCHAR(50) REFERENCES groups(name),
    roll_number  CHAR(8)     REFERENCES student(roll_number),
    PRIMARY KEY(group_name, roll_number)
);

CREATE TABLE IF NOT EXISTS group_social (
    group_name     VARCHAR(50) REFERENCES groups(name),
    platform_type  VARCHAR(15),
    link           VARCHAR     NOT NULL,
    PRIMARY KEY (group_name, platform_type)
);

CREATE OR REPLACE VIEW group_discord_user AS
    SELECT
        s.batch,
        s.discord_id,
        g.name, g.alias,
        gd.server_id, gs.link,
        gd.freshman_role, gd.sophomore_role, gd.junior_role, gd.senior_role,
        gd.guest_role
    FROM
        group_member AS gm
    JOIN groups AS g ON g.name = gm.group_name
    JOIN group_discord gd ON gd.group_name = gm.group_name
    JOIN group_social AS gs ON gs.group_name = gm.group_name AND gs.platform_type = 'discord'
    JOIN student AS s ON s.roll_number = gm.roll_number
    WHERE
        s.discord_id IS NOT NULL;
