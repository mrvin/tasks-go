CREATE TABLE IF NOT EXISTS buildings (
	id bigserial PRIMARY KEY,
	name TEXT NOT NULL,
	city TEXT NOT NULL,
	year SMALLINT NOT NULL,
	number_floors SMALLINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_city ON buildings(city);
CREATE INDEX IF NOT EXISTS idx_year ON buildings(year);
CREATE INDEX IF NOT EXISTS idx_number_floors ON buildings(number_floors);
