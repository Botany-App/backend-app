CREATE TABLE plants (
    ID UUID PRIMARY KEY,
    name_plant VARCHAR(50) NOT NULL,
    description_plant VARCHAR(100) NOT NULL,
    location_plant VARCHAR(50) NOT NULL,
    planting_time TIMESTAMP NOT NULL,
    irrigation_week NUMERIC NOT NULL,
    sun_exposure NUMERIC NOT NULL,
    fertilization_week NUMERIC NOT NULL,
    user_id UUID NOT NULL,
    species_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_plant FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE,
    CONSTRAINT fk_species_plant FOREIGN KEY (species_id) REFERENCES species(ID) ON DELETE CASCADE
);