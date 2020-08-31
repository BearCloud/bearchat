  
CREATE DATABASE auth;

USE auth;

CREATE TABLE users (
    email VARCHAR(320),
    hashedPassword TEXT,
    verified boolean,
    resetToken TEXT,
    verifiedToken TEXT,
    userId VARCHAR(128)
);