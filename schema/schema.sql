
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    email VARCHAR(100) UNIQUE,
    password VARCHAR(100) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    address VARCHAR(50) NOT NULL,
    profile_picture VARCHAR(100),
    email_verified  BOOLEAN,
	role   VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create product category table
CREATE TABLE IF NOT EXISTS product_category (
    category_id INT AUTO_INCREMENT PRIMARY KEY,
    parent_category INT,
    name VARCHAR(50),
    FOREIGN KEY (parent_category) REFERENCES product_category(category_id)
);

-- Create product table
CREATE TABLE IF NOT EXISTS product (
    product_id INT AUTO_INCREMENT PRIMARY KEY,
    category_id INT,
    product_name VARCHAR(50),
    description VARCHAR(100),
    FOREIGN KEY (category_id) REFERENCES product_category(category_id)
);