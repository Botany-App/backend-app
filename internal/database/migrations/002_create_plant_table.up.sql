CREATE TABLE IF NOT EXISTS plants (
			ID UUID PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			species VARCHAR(255),
			location_plant VARCHAR(255),
			planting_time TIMESTAMPTZ NOT NULL,
			water_times_per_week NUMERIC NOT NULL,
			sun_exposure_time_per_day NUMERIC NOT NULL,
			fertilizing_times_per_week NUMERIC NOT NULL,
			user_id UUID NOT NULL, 
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT fk_user_plant FOREIGN KEY(user_id) REFERENCES users(ID) ON DELETE CASCADE
		)