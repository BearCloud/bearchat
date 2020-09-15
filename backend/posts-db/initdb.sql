CREATE DATABASE postsDB;

USE postsDB;

CREATE TABLE posts (
    content VARCHAR(255),
    postID VARCHAR(36) PRIMARY KEY,
    authorID VARCHAR(36),
    postTime DATETIME
);
