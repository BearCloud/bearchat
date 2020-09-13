CREATE DATABASE postsDB;

USE postsDB;

CREATE TABLE posts (
    content VARCHAR(255),
    uuid VARCHAR(128),
    privacy INT,
    postTime DATETIME
);
