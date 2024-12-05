CREATE TABLE IF NOT EXISTS events (
    id SERIAL NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    date date NOT NULL,
    time time NOT NULL,
    venue_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT events_pkey PRIMARY KEY (id),
    CONSTRAINT events_venue_fk FOREIGN KEY (venue_id) REFERENCES venues (id)
    );

CREATE INDEX idx_events_deleted_at
    ON events USING btree
    (deleted_at ASC NULLS LAST);
