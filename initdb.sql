DROP DATABASE arturdb;

CREATE DATABASE arturdb; 

\c arturdb;

CREATE TABLE flashcards (
    id       BIGSERIAL PRIMARY KEY,
    question CHAR(512) NOT NULL,
    answer   CHAR(512) NOT NULL,
    date     TIMESTAMP
);