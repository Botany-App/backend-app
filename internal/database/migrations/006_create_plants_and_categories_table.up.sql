CREATE TABLE plants_and_categories (
    ID UUID PRIMARY KEY,
    plant_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_plant_category FOREIGN KEY (plant_id) REFERENCES plants(ID) ON DELETE CASCADE,
    CONSTRAINT fk_category_plant FOREIGN KEY (category_id) REFERENCES categories_plants(ID) ON DELETE CASCADE
);