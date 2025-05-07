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
CREATE INDEX IF NOT EXISTS idx_age ON persons(age);
CREATE INDEX IF NOT EXISTS idx_gender ON persons(gender);
CREATE INDEX IF NOT EXISTS idx_country_id ON persons(country_id);
