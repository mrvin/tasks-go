CREATE TABLE IF NOT EXISTS schedules (
	id bigserial PRIMARY KEY,
	name_medicine text,
	num_per_day smallint,
	times bigint[],
	all_life boolean,
	begin_date date,
	end_date date,
	user_id uuid
);
