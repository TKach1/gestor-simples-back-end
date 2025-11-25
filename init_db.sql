-- PostgreSQL Database Modeling Script for Gestor Simples

-- Table: Users
-- Stores information about users (administrators and sellers).
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL -- e.g., 'admin' or 'vendedor'
);

-- Table: Products
-- Stores information about products in stock.
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    price REAL NOT NULL DEFAULT 0.0
);

-- Table: Sales
-- Records all sales made in the system.
CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: Sales_Items
-- Stores the items within a sale.
CREATE TABLE IF NOT EXISTS sales_items (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER NOT NULL REFERENCES sales(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL
);

-- Optional: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_sales_user_id ON sales (user_id);
CREATE INDEX IF NOT EXISTS idx_sales_items_sale_id ON sales_items (sale_id);
CREATE INDEX IF NOT EXISTS idx_sales_items_product_id ON sales_items (product_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- Optional: Add a few initial users and products for testing
-- You might want to hash the password for 'admin' user with your application's hashing logic
-- INSERT INTO users (name, username, password_hash, role) VALUES ('Admin User', 'admin', 'your_hashed_admin_password', 'admin');
-- INSERT INTO users (name, username, password_hash, role) VALUES ('Seller One', 'seller1', 'your_hashed_seller1_password', 'vendedor');

-- INSERT INTO products (name, description, quantity, price) VALUES ('Product A', 'Description for Product A', 100, 19.99);
-- INSERT INTO products (name, description, quantity, price) VALUES ('Product B', 'Description for Product B', 50, 49.99);
