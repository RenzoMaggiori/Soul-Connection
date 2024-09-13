CREATE TABLE IF NOT EXISTS "employee" (
    id SERIAL PRIMARY KEY,
    soul_connection_id INT UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    birth_date VARCHAR(255) NOT NULL,
    gender VARCHAR(255) NOT NULL,
    work VARCHAR(255) NOT NULL,
    image_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "event" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date VARCHAR(255) NOT NULL,
    max_participants INT NOT NULL,
    location_x VARCHAR(255) NOT NULL,
    location_y VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    employee_id INT REFERENCES employee(id) ON DELETE CASCADE,
    CONSTRAINT unique_event UNIQUE (name, date, location_x, location_y)
);

CREATE TABLE IF NOT EXISTS "customer" (
    id SERIAL PRIMARY KEY,
    soul_connection_id INT UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    birth_date VARCHAR(255) NOT NULL,
    gender VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    astrological_sign VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    image_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    employee_id INT REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "payment" (
    id SERIAL PRIMARY KEY,
    soul_connection_id INT UNIQUE,
    date VARCHAR(255) NOT NULL,
    payment_method VARCHAR(255) NOT NULL,
    amount FLOAT(25) NOT NULL,
    comment VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    customer_id INT REFERENCES customer(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "encounter" (
    id SERIAL PRIMARY KEY,
    date VARCHAR(255) NOT NULL,
    rating INT NOT NULL,
    comment VARCHAR(255) NOT NULL,
    source VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    customer_id INT REFERENCES customer(id) ON DELETE CASCADE,
    CONSTRAINT unique_encounter UNIQUE (date, comment, source, customer_id)
);

CREATE TABLE IF NOT EXISTS "clothe" (
    id SERIAL PRIMARY KEY,
    soul_connection_id INT UNIQUE,
    type VARCHAR(255) NOT NULL,
    image_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    customer_id INT REFERENCES customer(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "tip" (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    tip VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_tip UNIQUE (title, tip)
);

