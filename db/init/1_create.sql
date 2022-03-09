CREATE DATABASE sample_blog;
use sample_blog;
CREATE TABLE entries (
    id INT(11) AUTO_INCREMENT NOT NULL,
    title VARCHAR(255),
    description VARCHAR(255),
    body VARCHAR(255),
    PRIMARY KEY (id)
);
