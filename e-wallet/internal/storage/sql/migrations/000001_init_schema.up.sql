CREATE TABLE IF NOT EXISTS wallets (
	id UUID DEFAULT gen_random_uuid (),
	balance double precision,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transactions (
	time timestamptz,
	from_wallet_id UUID references wallets(id),
	to_wallet_id UUID references wallets(id),
	amount double precision
);
