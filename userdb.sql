
 use userdb;
 CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    password VARCHAR(255),
    role ENUM('admin', 'courier', 'user') DEFAULT 'user'
);
CREATE TABLE orders (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    pickup_location VARCHAR(255),
    dropoff_location VARCHAR(255),
    package_details TEXT,
    delivery_time DATETIME,
    status ENUM('Pending', 'Picked Up', 'In Transit', 'Delivered') DEFAULT 'Pending',
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);
CREATE TABLE couriers (
    courier_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    password VARCHAR(255)
);
CREATE TABLE assigned_orders (
    assignment_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT,
    courier_id INT,
    status ENUM('Assigned', 'Accepted', 'Declined') DEFAULT 'Assigned',
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (courier_id) REFERENCES couriers(courier_id)
);
