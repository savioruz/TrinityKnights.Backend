CREATE TABLE IF NOT EXISTS venues (
    id SERIAL NOT NULL,
    name varchar(255) NOT NULL,
    address text,
    capacity integer,
    city varchar(255),
    state varchar(255),
    zip_code varchar(20),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT venues_pkey PRIMARY KEY (id)
    );

CREATE INDEX idx_venues_deleted_at
    ON venues USING btree
    (deleted_at ASC NULLS LAST);
