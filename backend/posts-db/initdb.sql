CREATE DATABASE postsDB;

USE postsDB;

CREATE TABLE posts (
    content VARCHAR(255),
    postID VARCHAR(128) PRIMARY KEY,
    uuid VARCHAR(128),
    postTime DATETIME
);
