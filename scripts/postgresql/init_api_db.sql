DROP DATABASE IF EXISTS web_forum_db;
CREATE DATABASE web_forum_db
    WITH OWNER postgres
    ENCODING 'utf8';
\connect web_forum_db;

CREATE EXTENSION IF NOT EXISTS citext;

-- User profile
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    username CITEXT NOT NULL,
    avatar TEXT NOT NULL DEFAULT '',
    email CITEXT NOT NULL,
    password BYTEA NOT NULL,
    date_created  TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    admin_level INTEGER DEFAULT 1,

    CONSTRAINT email_unique UNIQUE (email),
    CONSTRAINT username_unique UNIQUE (username)
);