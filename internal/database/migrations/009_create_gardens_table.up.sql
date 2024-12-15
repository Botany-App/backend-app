CREATE TABLE gardens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    garden_name VARCHAR(50) NOT NULL,
    garden_description TEXT NOT NULL,
    garden_location VARCHAR(50) NOT NULL,
    total_area NUMERIC(10, 3) NOT NULL,
    currenting_height NUMERIC(10, 2) NOT NULL,
    currenting_width NUMERIC(10, 2) NOT NULL,
    planting_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_irrigation TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_fertilization TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    irrigation_week NUMERIC NOT NULL,
    sun_exposure NUMERIC NOT NULL,
    fertilization_week NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_garden FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE garden_plant(
    id UUID PRIMARY KEY,
    garden_id UUID NOT NULL,
    plant_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_garden_plant FOREIGN KEY (garden_id) REFERENCES gardens(id) ON DELETE CASCADE,
    CONSTRAINT fk_plant_garden FOREIGN KEY (plant_id) REFERENCES plants(id) ON DELETE CASCADE
)