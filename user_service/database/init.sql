CREATE DATABASE  user_service;

-- Создание последовательности
CREATE SEQUENCE users_id_seq START 1 INCREMENT 1;

-- Создание таблицы с явным указанием последовательности
CREATE TABLE user (
    id INTEGER PRIMARY KEY DEFAULT nextval('users_id_seq'),
    name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

-- Привязка последовательности к таблице (опционально)
ALTER SEQUENCE users_id_seq OWNED BY users.id;

ALTER TABLE public."user" ADD email varchar NOT NULL;
ALTER TABLE public."user" ADD CONSTRAINT user_unique UNIQUE (email);
