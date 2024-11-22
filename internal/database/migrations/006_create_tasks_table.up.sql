CREATE TYPE task_status_enum AS ENUM ('pending', 'in_progress', 'completed');

CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    task_name VARCHAR(50) NOT NULL,
    task_description VARCHAR(100) NOT NULL,
    date_task TIMESTAMP NOT NULL,
    urgency_level  NUMERIC CHECK (urgency_level BETWEEN 1 AND 5) NOT NULL,
    task_status task_status_enum NOT NULL DEFAULT 'pending',
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_task FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

