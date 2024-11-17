CREATE TABLE task_plants (
    ID UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    plant_id UUID NOT NULL,
    CONSTRAINT fk_task_plant FOREIGN KEY (task_id) REFERENCES tasks(ID) ON DELETE CASCADE,
    CONSTRAINT fk_plant_task FOREIGN KEY (plant_id) REFERENCES plants(ID) ON DELETE CASCADE
);

CREATE TABLE task_gardens (
    ID UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    garden_id UUID NOT NULL,
    CONSTRAINT fk_task_garden FOREIGN KEY (task_id) REFERENCES tasks(ID) ON DELETE CASCADE,
    CONSTRAINT fk_garden_task FOREIGN KEY (garden_id) REFERENCES gardens(ID) ON DELETE CASCADE
);


CREATE TRIGGER trigger_update_timestamp_task_gardens
BEFORE UPDATE ON task_gardens
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


CREATE TRIGGER trigger_update_timestamp_task_plants
BEFORE UPDATE ON task_plants
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

