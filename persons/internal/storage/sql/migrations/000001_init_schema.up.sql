CREATE TABLE IF NOT EXISTS persons (
	id serial primary key,
	name text,
	surname text,
	patronymic text,
	age int,
	gender text,
	country_id text
);