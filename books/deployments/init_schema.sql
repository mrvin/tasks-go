CREATE TABLE IF NOT EXISTS books (
	id serial primary key,
	title varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS authors (
	id serial primary key,
	name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS book_author (
	id_book integer references books(id) on delete cascade,
	id_author integer references authors(id)
);

