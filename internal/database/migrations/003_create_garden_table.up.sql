CREATE TABLE IF NOT EXISTS gardens(
  	ID UUID PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		location_garden VARCHAR(255),
		water_times_per_week NUMERIC NOT NULL,
		sun_exposure_time_per_day NUMERIC NOT NULL,
		fertilizing_times_per_week NUMERIC NOT NULL,
		user_id UUID NOT NULL, 
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		CONSTRAINT fk_user_garden FOREIGN KEY(user_id) REFERENCES users(ID) ON DELETE CASCADE
);

-- Tabela gardenPlant
CREATE TABLE IF NOT EXISTS gardenPlant (
		garden_id UUID NOT NULL,
		plant_id UUID NOT NULL,
		CONSTRAINT fk_garden_gardenPlant FOREIGN KEY(garden_id) REFERENCES gardens(ID) ON DELETE CASCADE,
		CONSTRAINT fk_plant_gardenPlant FOREIGN KEY(plant_id) REFERENCES plants(ID) ON DELETE CASCADE
);