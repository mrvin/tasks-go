CREATE TABLE IF NOT EXISTS persons (
	id SERIAL NOT NULL UNIQUE PRIMARY KEY,
	name TEXT,
	surname TEXT,
	patronymic TEXT,
	age INT,
	gender TEXT,
	country_id TEXT
);