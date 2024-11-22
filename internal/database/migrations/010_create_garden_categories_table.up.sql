CREATE TABLE garden_categories (
    id UUID PRIMARY KEY,
    garden_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_garden_category FOREIGN KEY (garden_id) REFERENCES gardens(id) ON DELETE CASCADE,
    CONSTRAINT fk_category_garden FOREIGN KEY (category_id) REFERENCES categories_plants(id) ON DELETE CASCADE
);