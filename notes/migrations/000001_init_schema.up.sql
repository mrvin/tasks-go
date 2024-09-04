CREATE TABLE IF NOT EXISTS users (
	name TEXT NOT NULL UNIQUE,
	hash_password TEXT NOT NULL,
	role TEXT NOT NULL,
	PRIMARY KEY (name)
);
CREATE INDEX IF NOT EXISTS idx_name ON users(name);

CREATE TABLE IF NOT EXISTS notes (
	id serial primary key,
	title text,
	description text,
	user_name TEXT references users(name) on delete cascade
);
