-- Connect to the newly created database
\c stocksdb

-- Create the table
CREATE TABLE stocks (
    stockid SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    company VARCHAR(255)
);
