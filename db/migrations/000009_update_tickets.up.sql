ALTER TABLE tickets ALTER COLUMN seat_number TYPE varchar(10) USING seat_number::varchar(10);
