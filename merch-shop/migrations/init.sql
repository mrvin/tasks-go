CREATE TABLE IF NOT EXISTS accounts (
	name TEXT NOT NULL UNIQUE PRIMARY KEY,
	hash_password TEXT NOT NULL,
	balance NUMERIC(20,0) NOT NULL CHECK (balance >= 0 AND balance <= 18446744073709551615) DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_name ON accounts(name);

CREATE TABLE IF NOT EXISTS transactions (
	time timestamptz,
	from_name TEXT references accounts(name),
	to_name TEXT references accounts(name),
	amount NUMERIC(20,0) NOT NULL CHECK (amount >= 0 AND amount <= 18446744073709551615)
);

CREATE TABLE IF NOT EXISTS products (
	name TEXT NOT NULL UNIQUE PRIMARY KEY,
	price NUMERIC(20,0) NOT NULL CHECK (price >= 0 AND price <= 18446744073709551615)
);

CREATE TABLE IF NOT EXISTS orders (
	user_name TEXT references accounts(name),
	product_name TEXT references products(name)
);

INSERT INTO products (name, price)
	VALUES ('t-shirt', 80),
		('cup', 20),
		('book', 50),
		('pen', 10),
		('powerbank', 200),
		('hoody', 300),
		('umbrella', 200),
		('socks', 10),
		('wallet', 50),
		('pink-hoody', 500);
