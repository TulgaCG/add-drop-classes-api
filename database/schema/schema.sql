CREATE TYPE USER_ROLE_VALUE AS ENUM ('admin', 'teacher', 'student');

CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL,
    username        TEXT NOT NULL,
    password        TEXT NOT NULL,
    token           TEXT,
    token_expire_at TIMESTAMP,

    PRIMARY KEY(id),
    UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS roles (
    id   BIGSERIAL,
    role USER_ROLE_VALUE NOT NULL,

    PRIMARY KEY(id),
    UNIQUE (role)
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGSERIAL NOT NULL,
    role_id BIGSERIAL NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    UNIQUE(user_id, role_id)
);

CREATE TABLE IF NOT EXISTS lectures (
    id     BIGSERIAL,
    name   TEXT NOT NULL,
    code   TEXT NOT NULL,
    credit SMALLINT NOT NULL,
    type   SMALLINT,

    PRIMARY KEY(id),
    UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS user_lectures (
    id         BIGSERIAL,
    user_id    BIGSERIAL NOT NULL,
    lecture_id BIGSERIAL NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE,
    UNIQUE(user_id, lecture_id)
);

INSERT INTO roles(id, role)
VALUES
    (1, 'admin'),
    (2, 'teacher'),
    (3, 'student')
ON CONFLICT DO NOTHING;
