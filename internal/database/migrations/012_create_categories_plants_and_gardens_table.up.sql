CREATE TABLE categories_plants_and_gardens (
    garden_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_garden_category FOREIGN KEY (garden_id) REFERENCES gardens(ID) ON DELETE CASCADE,
    CONSTRAINT fk_category_garden FOREIGN KEY (category_id) REFERENCES categories_plants(ID) ON DELETE CASCADE
);