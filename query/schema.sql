CREATE TABLE IF NOT EXISTS hostel (
    number       varchar(4) primary key check(number like '%H_%'),
    name         text       unique,
    warden_name  text       unique
);

CREATE TABLE IF NOT EXISTS student (
    roll_number    int         primary key,
    section        char(4)     not null check(section like '__-_'),
    sub_section    char(5)     not null check(sub_section like '__-__'),
    name           varchar     not null,
    gender         char(1)     check(gender='M' or gender='F' or gender='O'),
    mobile         varchar(14) unique,
    birthday       date,
    email          text        not null check(email like '%___@nitkkr.ac.in'),
    batch          smallint    not null check(batch >= 0),
    hostel_number  varchar(4)  references hostel(number),
    room_number    varchar(6)  check(room_number like '%_-___'),
    discord_uid    bigint      unique,
    verified       boolean     not null default false
);

CREATE TABLE IF NOT EXISTS course (
    code        char(7) NOT NULL,
    title       varchar(32) NOT NULL,
    branch      char(3) NOT NULL,
    semester    smallint NOT NULL,
    credits     smallint[4] NOT NULL,
    prereq      char(7)[],
    type        char(3) NOT NULL,
    objectives  text NOT NULL,
    content     text NOT NULL,
    books       text NOT NULL,
    outcomes    text NOT NULL,
    PRIMARY KEY (code, branch)
);

CREATE TABLE IF NOT EXISTS groups (
    name         text       primary key,
    alias        text       unique,
    branch       varchar(5) unique,
    kind         text       NOT NULL check(kind in ('cultural club', 'technical club', 'technical society')),
    description  text
);

CREATE TABLE IF NOT EXISTS group_faculty (
    group_name  text   references groups(name),
    name        text,
    mobile      bigint unique NOT NULL,
    primary key (group_name, name)
);

CREATE TABLE IF NOT EXISTS group_discord (
    name            text        references groups(name),
    id              bigint      unique,
    invite          varchar(10) unique,
    fresher_role    bigint      unique,
    sophomore_role  bigint      unique,
    junior_role     bigint      unique,
    senior_role     bigint      unique,
    guest_role      bigint      unique
);

CREATE TABLE IF NOT EXISTS group_social (
    name  text references groups(name),
    type  varchar(15),
    link  text NOT NULL,
    primary key (name, type)
);

CREATE TABLE IF NOT EXISTS group_admin (
    group_name   text references groups(name),
    position     varchar(20) NOT NULL,
    roll_number  int  references student(roll_number),
    primary key (group_name, roll_number)
);

CREATE TABLE IF NOT EXISTS group_member (
    roll_number  int  not null references student(roll_number),
    group_name   text not null references groups(name),
    primary key(roll_number, group_name)
);

CREATE OR REPLACE VIEW group_discord_user AS
    SELECT student.batch, student.discord_uid,
        c.name, c.alias, d.id, d.invite,
        d.fresher_role, d.sophomore_role, d.junior_role, d.senior_role,
        d.guest_role
        FROM group_member m
        JOIN groups c ON c.name = m.group_name
        JOIN group_discord d ON d.name = m.group_name
        JOIN student ON student.roll_number = m.roll_number
        WHERE student.discord_uid IS NOT NULL;
