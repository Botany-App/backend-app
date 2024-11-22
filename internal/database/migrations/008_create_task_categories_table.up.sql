CREATE TABLE task_categories(
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_task_category FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_category_task FOREIGN KEY (category_id) REFERENCES categories_tasks(id) ON DELETE CASCADE
);
