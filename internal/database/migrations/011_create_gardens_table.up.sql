CREATE TABLE gardens (
    ID UUID PRIMARY KEY,
    name_garden VARCHAR(50) NOT NULL,
    description_garden VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    location_garden VARCHAR(50) NOT NULL,
    irrigation_week NUMERIC NOT NULL,
    sun_exposure NUMERIC NOT NULL,
    fertilization_week NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_garden FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);