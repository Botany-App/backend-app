CREATE TABLE species (
    id UUID PRIMARY KEY,
    common_name VARCHAR(50) NOT NULL,
    specie_description VARCHAR(100) NOT NULL,
    scientific_name varchar(50) NOT NULL,
    botanical_family varchar(50) NOT NULL,
    growth_type varchar(50) NOT NULL,
    ideal_temperature NUMERIC NOT NULL,
    ideal_climate varchar(50) NOT NULL,
    life_cycle VARCHAR(50) NOT NULL,
    planting_season VARCHAR(50) NOT NULL,
    harvest_time INTEGER NOT NULL,
    average_height NUMERIC(10, 2) NOT NULL,
    average_width NUMERIC(10, 2) NOT NULL,
    irrigation_weight NUMERIC(5, 2) NOT NULL,
    fertilization_weight NUMERIC(5, 2) NOT NULL,
    sun_weight NUMERIC(5, 2) NOT NULL,
    image_url VARCHAR(300) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

