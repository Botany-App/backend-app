CREATE TABLE categories_tasks_and_tasks (
    ID UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_task_category FOREIGN KEY (task_id) REFERENCES tasks(ID) ON DELETE CASCADE,
    CONSTRAINT fk_category_task FOREIGN KEY (category_id) REFERENCES categories_tasks(ID) ON DELETE CASCADE
);
