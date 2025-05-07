CREATE TYPE gender_type AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS persons (
	id SERIAL NOT NULL UNIQUE PRIMARY KEY,
	name TEXT,
	surname TEXT,
	patronymic TEXT,
	age INT,
	gender gender_type,
	country_id TEXT
);