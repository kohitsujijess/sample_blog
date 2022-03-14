DROP DATABASE IF EXISTS sample_blog_test;
CREATE DATABASE IF NOT EXISTS sample_blog_test;
USE sample_blog_test;

DROP TABLE IF EXISTS entries;
CREATE TABLE IF NOT EXISTS entries (
  id VARCHAR(191) PRIMARY KEY, 
  title LONGTEXT,
  body LONGTEXT,
  created_at DATETIME(3),
  updated_at DATETIME(3)
);
