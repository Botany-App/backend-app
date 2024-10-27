CREATE TABLE history_plants (
    ID UUID PRIMARY KEY,
    plant_id UUID NOT NULL,
    location_garden VARCHAR(50) NOT NULL,
    irrigation_week NUMERIC NOT NULL,
    sun_exposure NUMERIC NOT NULL,
    fertilization_week NUMERIC NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_plant_log FOREIGN KEY (plant_id) REFERENCES plants(ID) ON DELETE CASCADE,
    CONSTRAINT fk_user_log FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);

CREATE TABLE history_gardens(
    ID UUID PRIMARY KEY,
    garden_id UUID NOT NULL,
    location_garden VARCHAR(50) NOT NULL,
    irrigation_week NUMERIC NOT NULL,
    sun_exposure NUMERIC NOT NULL,
    fertilization_week NUMERIC NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_garden_log FOREIGN KEY (garden_id) REFERENCES gardens(ID) ON DELETE CASCADE,
    CONSTRAINT fk_user_log FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
)