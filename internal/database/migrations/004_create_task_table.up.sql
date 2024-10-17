CREATE TYPE task_status_enum AS ENUM ('pending', 'in_progress', 'completed');
CREATE TABLE IF NOT EXISTS task (
    ID UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(255),
    task_description TEXT,
    task_status task_status_enum NOT NULL DEFAULT 'pending', 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    due_date TIMESTAMPTZ,
    last_updated TIMESTAMPTZ DEFAULT NOW(),
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_task FOREIGN KEY(user_id) REFERENCES users(ID) ON DELETE CASCADE
);


-- Tabela taskGarden
CREATE TABLE IF NOT EXISTS taskGarden (
    task_id UUID NOT NULL,
    garden_id UUID NOT NULL,
    CONSTRAINT fk_task_taskGarden FOREIGN KEY(task_id) REFERENCES task(ID) ON DELETE CASCADE,
    CONSTRAINT fk_garden_taskGarden FOREIGN KEY(garden_id) REFERENCES gardens(ID) ON DELETE CASCADE
);

-- Tabela taskPlant
CREATE TABLE IF NOT EXISTS taskPlant (
    task_id UUID NOT NULL,
    plant_id UUID NOT NULL,
    CONSTRAINT fk_task_taskPlant FOREIGN KEY(task_id) REFERENCES task(ID) ON DELETE CASCADE,
    CONSTRAINT fk_plant_taskPlant FOREIGN KEY(plant_id) REFERENCES plants(ID) ON DELETE CASCADE
);




