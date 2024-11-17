CREATE TYPE task_status_enum AS ENUM ('pending', 'in_progress', 'completed');

CREATE TABLE tasks (
    ID UUID PRIMARY KEY,
    name_task VARCHAR(50) NOT NULL,
    description_task VARCHAR(100) NOT NULL,
    date_task TIMESTAMP NOT NULL,
    task_status task_status_enum NOT NULL DEFAULT 'pending',
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_task FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);

