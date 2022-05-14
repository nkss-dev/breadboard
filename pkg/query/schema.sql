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
