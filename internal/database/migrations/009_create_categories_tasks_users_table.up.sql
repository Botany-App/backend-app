CREATE TABLE categories_tasks_users(
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_user_category FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE,
    CONSTRAINT fk_category_user FOREIGN KEY (category_id) REFERENCES categories_tasks(ID) ON DELETE CASCADE
);