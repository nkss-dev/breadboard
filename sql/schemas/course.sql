CREATE TABLE IF NOT EXISTS book (
    book_name     VARCHAR   PRIMARY KEY,
    author_names  VARCHAR[] NOT NULL,
    publisher     VARCHAR   NOT NULL,
    edition       SMALLINT  NOT NULL,
    url           VARCHAR
);

CREATE TABLE IF NOT EXISTS course (
    code        VARCHAR(7)  PRIMARY KEY,  /* code can vary from 6-7 characters */
    title       VARCHAR(64) NOT NULL,
    prereq      VARCHAR(7)[],             /* prereq. of this course in course code format */ 
    kind        CHAR(3)     NOT NULL,
    objectives  VARCHAR[]   NOT NULL,
    content     VARCHAR     NOT NULL,
    book_names  VARCHAR[]   NOT NULL,
    outcomes    VARCHAR[]   NOT NULL
);

CREATE TABLE IF NOT EXISTS branch_specifics (
    code      VARCHAR(7)  REFERENCES course(code),
    branch    VARCHAR(3)  NOT NULL,  /* branch can very from 2-3 characters */
    semester  SMALLINT    NOT NULL CHECK(semester IN (1, 2, 3, 4, 5, 6, 7, 8)),
    credits   SMALLINT[4] NOT NULL,  /* lecture, tutorial, practical and total credits respectively */
    PRIMARY KEY (code, branch)
);
