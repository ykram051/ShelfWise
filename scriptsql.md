CREATE TABLE IF NOT EXISTS authors (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    bio TEXT
);

CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INT NOT NULL REFERENCES authors(id),
    genres TEXT[],
    published_at TIMESTAMP NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    stock INT NOT NULL
);

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    street TEXT,
    city TEXT,
    state TEXT,
    postal_code TEXT,
    country TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    total_price DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    book_id INT NOT NULL REFERENCES books(id),
    quantity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS sales_reports (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT current_timestamp,
    total_revenue DOUBLE PRECISION NOT NULL,
    total_orders INT NOT NULL
);
