CREATE DATABASE IF NOT EXISTS sample_blog;
USE sample_blog;
CREATE TABLE IF NOT EXISTS entries (
    id INT(11) AUTO_INCREMENT NOT NULL,
    uuid VARCHAR(255) unique,
    title VARCHAR(255),
    description VARCHAR(255),
    body text,
    created_at datetime  default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
);
