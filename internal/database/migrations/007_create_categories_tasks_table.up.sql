CREATE TABLE categories_tasks (
    id UUID PRIMARY KEY,
    category_name VARCHAR(50) NOT NULL,
    category_description VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_category FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);