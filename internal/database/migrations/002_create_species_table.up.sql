CREATE TABLE species (
    ID UUID PRIMARY KEY,
    name_species VARCHAR(50) NOT NULL,
    description_species VARCHAR(100) NOT NULL,
    fertilization_weight NUMERIC NOT NULL,
    sun_weight NUMERIC NOT NULL,
    irrigation_weight NUMERIC NOT NULL,
    time_to_harvest NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE EXTENSION IF NOT EXISTS "pgcrypto";


INSERT INTO species (
    ID, name_species, description_species, fertilization_weight, sun_weight, irrigation_weight, time_to_harvest, created_at, updated_at
) VALUES 
    (gen_random_uuid(), 'Laranja', 'Fruto cítrico rico em vitamina C', 1.5, 2.0, 1.8, 180, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'Milho', 'Cereal amplamente cultivado para alimento', 1.3, 2.5, 2.1, 120, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
