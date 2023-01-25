CREATE TABLE IF NOT EXISTS academic_announcement (
    date_of_creation  DATE    NOT NULL,
    title             VARCHAR NOT NULL,
    title_link        VARCHAR NOT NULL,
    subtitle          VARCHAR,
    subtitle_link     VARCHAR,
    kind              VARCHAR(20) NOT NULL
    CHECK (kind IN (
        'academic',
        'result',
        'exam'
    )),
    PRIMARY KEY (date_of_creation, title)
);
