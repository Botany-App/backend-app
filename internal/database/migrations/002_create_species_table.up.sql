CREATE TABLE species (
    ID UUID PRIMARY KEY,
    name_species VARCHAR(50) NOT NULL,
    description_species VARCHAR(100) NOT NULL,
    fertilization_weight NUMERIC NOT NULL,
    sun_weight NUMERIC NOT NULL,
    irrigation_weight NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);