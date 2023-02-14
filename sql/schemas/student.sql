CREATE TABLE IF NOT EXISTS hostel (
    id     VARCHAR(4)  PRIMARY KEY CHECK(id LIKE '%H_%'),
    name   VARCHAR(40) UNIQUE NOT NULL,
    email  VARCHAR(40) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS warden (
    name       VARCHAR(50) PRIMARY KEY,
    mobile     VARCHAR(14) UNIQUE NOT NULL,
    hostel_id  VARCHAR(4)  NOT NULL REFERENCES hostel(id)
);

CREATE TABLE IF NOT EXISTS student (
    roll_number  CHAR(8)     PRIMARY KEY,
    section      VARCHAR(6)  NOT NULL CHECK(section SIMILAR TO '__-_[1-9][0-9]?'),
    name         VARCHAR(50) NOT NULL,
    gender       CHAR(1)     CHECK(gender IN ('M', 'F', 'O')),
    mobile       VARCHAR(14) UNIQUE,
    birth_date   DATE,
    email        VARCHAR(40) NOT NULL CHECK(email LIKE '%___@nitkkr.ac.in'),
    batch        SMALLINT    NOT NULL CHECK(batch >= 0),
    hostel_id    VARCHAR(4)  NOT NULL REFERENCES hostel(id),
    room_id      VARCHAR(6)  CHECK(room_id LIKE '%_-___'),
    discord_id   BIGINT      UNIQUE,
    is_verified  BOOLEAN     DEFAULT false NOT NULL
    clubs        JSONB       DEFAULT '{}'::JSONB NOT NULL,
);
