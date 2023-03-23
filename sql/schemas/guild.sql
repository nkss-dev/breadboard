CREATE TABLE IF NOT EXISTS guild (
    id          BIGINT       PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    locale      VARCHAR(5)   DEFAULT 'en-GB' NOT NULL,
    bot_role    BIGINT       NOT NULL,
    mute_role   BIGINT,
    edit_log    BIGINT,
    delete_log  BIGINT
);

CREATE TABLE IF NOT EXISTS guild_role (
    guild_id  BIGINT      REFERENCES guild(id),
    field     VARCHAR(20) NOT NULL,
    value     VARCHAR(50),
    role_ids  BIGINT[],
    PRIMARY KEY (guild_id, value)
);

CREATE TABLE IF NOT EXISTS bot_prefix (
    guild_id  BIGINT NOT NULL REFERENCES guild(id),
    prefix    VARCHAR(5),
    PRIMARY KEY (guild_id, prefix)
);

CREATE TABLE IF NOT EXISTS join_role (
    guild_id  BIGINT NOT NULL REFERENCES guild(id),
    role_id   BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS mod_role (
    guild_id  BIGINT NOT NULL REFERENCES guild(id),
    role_id   BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS affiliated_guild (
    guild_id         BIGINT    PRIMARY KEY REFERENCES guild(id),
    batch            SMALLINT  UNIQUE NOT NULL CHECK(batch >= 0),
    info_channel     BIGINT    NOT NULL,
    command_channel  BIGINT    NOT NULL,
    guest_role       BIGINT    NOT NULL
);

CREATE TABLE IF NOT EXISTS event (
    guild_id    BIGINT       REFERENCES guild(id),
    event_type  VARCHAR(5)  NOT NULL,
    channel_id  BIGINT       NOT NULL,
    message     VARCHAR,
    PRIMARY KEY (guild_id, event_type),
    CONSTRAINT ck_event_type CHECK (
        event_type IN ('ban', 'kick', 'leave')
    )
);
