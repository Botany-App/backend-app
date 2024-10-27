CREATE TABLE categories_tasks (
    ID UUID PRIMARY KEY,
    name_category VARCHAR(50) NOT NULL,
    description_category VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);