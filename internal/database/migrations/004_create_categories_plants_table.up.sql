CREATE TABLE categories_plants (
    id UUID PRIMARY KEY,
    category_name VARCHAR(50) NOT NULL,
    category_description TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_category FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);