-- Create stock_entries table
CREATE TABLE IF NOT EXISTS stock_entries (
    id SERIAL PRIMARY KEY,
    hospital_id INTEGER NOT NULL,
    drug_name VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create dispensations table
CREATE TABLE IF NOT EXISTS dispensations (
    id SERIAL PRIMARY KEY,
    hospital_id INTEGER NOT NULL,
    drug_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indices for performance
CREATE INDEX IF NOT EXISTS idx_stock_hospital_drug ON stock_entries(hospital_id, drug_name);
CREATE INDEX IF NOT EXISTS idx_dispensations_hospital_drug ON dispensations(hospital_id, drug_name);
