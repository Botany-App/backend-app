CREATE TABLE task_plants (
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    plant_id UUID NOT NULL,
    CONSTRAINT fk_task_plant FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_plant_task FOREIGN KEY (plant_id) REFERENCES plants(id) ON DELETE CASCADE
);

CREATE TABLE task_gardens (
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    garden_id UUID NOT NULL,
    CONSTRAINT fk_task_garden FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_garden_task FOREIGN KEY (garden_id) REFERENCES gardens(id) ON DELETE CASCADE
);

