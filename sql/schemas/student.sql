CREATE TABLE IF NOT EXISTS hostel (
    id           VARCHAR(4)  PRIMARY KEY CHECK(id LIKE '%H_%'),
    hostel_name  VARCHAR(40) UNIQUE NOT NULL,
    warden_name  VARCHAR(50) UNIQUE
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
    hostel_id    VARCHAR(4)  REFERENCES hostel(id),
    room_id      VARCHAR(6)  CHECK(room_id LIKE '%_-___'),
    discord_id   BIGINT      UNIQUE,
    is_verified  BOOLEAN     DEFAULT false NOT NULL
);

-- CREATE TABLE IF NOT EXISTS ign (
--     discord_id     BIGINT PRIMARY KEY REFERENCES student(discord_id) ON UPDATE CASCADE,
--     app            TEXT,
--     kind           VARCHAR(20),
--     username       VARCHAR(30)
-- );

-- CREATE TABLE IF NOT EXISTS create_channel (
--     id            INT PRIMARY KEY,
--     create_text   INT,
--     create_voice  INT,
--     locked        INT
-- );

-- CREATE TABLE IF NOT EXISTS active_channel (
--     id              INT PRIMARY KEY,
--     master_channel  INT
-- );

-- CREATE TABLE IF NOT EXISTS voice (
--     voice_channel  INT UNIQUE,
--     owner          INT UNIQUE REFERENCES student(discord_id),
--     lock           INT CHECK (lock=1 or lock=0) default 0,
--     create_vc      INT CHECK (create_vc=1 or create_vc=0) default 0,
--     allow_text     INT CHECK (allow_text=1 or allow_text=0) default 1,
--     text_channel   INT
-- );
