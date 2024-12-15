CREATE TABLE history_plants (
    id UUID PRIMARY KEY DEFAULT,
    plant_id UUID NOT NULL,
    irrigation_week NUMERIC DEFAULT 0 NOT NULL,
    record_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    height NUMERIC(10, 2),
    width NUMERIC(10, 2),
    health_status VARCHAR(50) DEFAULT 'Saudável',
    irrigation BOOLEAN DEFAULT FALSE,
    fertilization BOOLEAN DEFAULT FALSE,
    sun_exposure NUMERIC NOT NULL,
    fertilization_week NUMERIC DEFAULT 0 NOT NULL,
    notes TEXT,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_plant_log FOREIGN KEY (plant_id) REFERENCES plants(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_log FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE history_gardens (
    id UUID PRIMARY KEY ,
    garden_id UUID NOT NULL,
    garden_location VARCHAR(50) NOT NULL,
    total_area NUMERIC(10, 3) NOT NULL,
    record_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    height NUMERIC(10, 2),
    width NUMERIC(10, 2),
    health_status VARCHAR(50) DEFAULT 'Saudável',
    irrigation BOOLEAN DEFAULT FALSE,
    fertilization BOOLEAN DEFAULT FALSE,
    irrigation_week NUMERIC DEFAULT 0 NOT NULL,
    sun_exposure NUMERIC NOT NULL,
    fertilization_week NUMERIC DEFAULT 0 NOT NULL,
    notes TEXT,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_garden_log FOREIGN KEY (garden_id) REFERENCES gardens(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_log FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_record_date_plants ON history_plants(record_date);
CREATE INDEX idx_record_date_gardens ON history_gardens(record_date);
