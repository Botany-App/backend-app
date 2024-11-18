CREATE TABLE categories_plants (
    ID UUID PRIMARY KEY,
    name_category VARCHAR(50) NOT NULL,
    description_category VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_category FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);