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

CREATE OR REPLACE PROCEDURE transfer(
	_from_wallet_id UUID,
	_to_wallet_id UUID,
	_amount double precision
)
LANGUAGE plpgsql
AS $$
DECLARE
balanceFrom double precision;
balanceTo double precision;
BEGIN
	SELECT balance INTO balanceFrom FROM wallets WHERE id = _from_wallet_id;
	IF NOT FOUND THEN
		RAISE EXCEPTION 'wallet % not found', _from_wallet_id;
	END IF;
	IF balanceFrom - _amount < 0 THEN
		RAISE EXCEPTION 'wallet % not enough funds', _from_wallet_id;
	END IF;

	SELECT balance INTO balanceTo FROM wallets WHERE id = _to_wallet_id;
	IF NOT FOUND THEN
		RAISE EXCEPTION 'wallet % not found', _to_wallet_id;
	END IF;

	UPDATE wallets SET balance = round(CAST(balance-_amount AS numeric), 2) WHERE id = _from_wallet_id;
	UPDATE wallets SET balance = round(CAST(balance+_amount AS numeric), 2) WHERE id = _to_wallet_id;

	INSERT INTO transactions (
		time,
		from_wallet_id,
		to_wallet_id,
		amount
	) VALUES (
		NOW(),
		_from_wallet_id,
		_to_wallet_id,
		_amount
	);

END;
$$;

CREATE OR REPLACE PROCEDURE withdraw(
	_from_wallet_id UUID,
	_amount double precision
)
LANGUAGE plpgsql
AS $$
DECLARE
balanceFrom double precision;
BEGIN
	SELECT balance INTO balanceFrom FROM wallets WHERE id = _from_wallet_id;
	IF NOT FOUND THEN
		RAISE EXCEPTION 'wallet % not found', _from_wallet_id;
	END IF;
	IF balanceFrom - _amount < 0 THEN
		RAISE EXCEPTION 'wallet % not enough funds', _from_wallet_id;
	END IF;


	UPDATE wallets SET balance = round(CAST(balance-_amount AS numeric), 2) WHERE id = _from_wallet_id;

	INSERT INTO transactions (
		time,
		from_wallet_id,
		to_wallet_id,
		amount
	) VALUES (
		NOW(),
		_from_wallet_id,
		NULL,
		_amount
	);

END;
$$;

CREATE OR REPLACE PROCEDURE deposit(
	_to_wallet_id UUID,
	_amount double precision
)
LANGUAGE plpgsql
AS $$
DECLARE
BEGIN
	UPDATE wallets SET balance = round(CAST(balance+_amount AS numeric), 2) WHERE id = _to_wallet_id;

	INSERT INTO transactions (
		time,
		from_wallet_id,
		to_wallet_id,
		amount
	) VALUES (
		NOW(),
		NULL,
		_to_wallet_id,
		_amount
	);

END;
$$;
