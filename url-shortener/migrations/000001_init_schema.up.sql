CREATE TABLE IF NOT EXISTS users (
	name TEXT NOT NULL UNIQUE PRIMARY KEY,
	hash_password TEXT NOT NULL,
	role TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_name ON users(name);

CREATE TABLE IF NOT EXISTS url(
	alias TEXT NOT NULL UNIQUE PRIMARY KEY,
	url TEXT NOT NULL,
	count NUMERIC(20,0) NOT NULL CHECK (count >= 0 AND count <= 18446744073709551615) DEFAULT 0,
	user_name TEXT references users(name) on delete cascade,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

CREATE OR REPLACE FUNCTION get_url(alias_for_url TEXT)
RETURNS TEXT
LANGUAGE plpgsql
as
$$
DECLARE
	full_url TEXT;
BEGIN
	SELECT url INTO full_url FROM url WHERE alias = alias_for_url;
	IF NOT FOUND THEN
		RAISE EXCEPTION 'alias % not found', alias_for_url;
	END IF;

	UPDATE url SET count = count+1 WHERE alias = alias_for_url;

	return full_url;
END;
$$;
