create database userdb;
use userdb;

CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    password VARCHAR(255),
    role ENUM('admin', 'courier', 'user') DEFAULT 'user'
);