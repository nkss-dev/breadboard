CREATE TABLE IF NOT EXISTS faculty (
    emp_id            INT         PRIMARY KEY,
    name              VARCHAR(50) NOT NULL,
    dp_url            VARCHAR     NOT NULL,
    designation       VARCHAR(20) NOT NULL,
    qualification     VARCHAR     NOT NULL,
    area_of_interest  VARCHAR[]   NOT NULL,
    experience        VARCHAR[]   NOT NULL,
    mobile            VARCHAR(14) NOT NULL,
    mobile_2          VARCHAR(14),
    email             VARCHAR(30) NOT NULL,
    department        VARCHAR[],
    is_regular        BOOLEAN DEFAULT true NOT NULL
);

CREATE TABLE IF NOT EXISTS hod (
    department  VARCHAR PRIMARY KEY,
    emp_id      INT     UNIQUE NOT NULL REFERENCES faculty(emp_id)
);
