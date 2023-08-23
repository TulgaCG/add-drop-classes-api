CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(64),
    token_expire_at DATETIME
);

CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY,
    role VARCHAR(255) NOT NULL,
    UNIQUE(id, role)
);

CREATE TABLE IF NOT EXISTS user_roles (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    UNIQUE(user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS lectures (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(6) UNIQUE NOT NULL,
    credit INTEGER NOT NULL,
    type INTEGER
);

CREATE TABLE IF NOT EXISTS user_lectures (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    lecture_id INTEGER NOT NULL,
    UNIQUE(user_id, lecture_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE
);