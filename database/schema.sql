CREATE TABLE IF NOT EXISTS users (
    id       UNSIGNED BIG INT,
    username VARCHAR(255),
    password VARCHAR(255),
    
    PRIMARY KEY (id),
    UNIQUE(username)
);