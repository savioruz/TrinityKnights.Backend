CREATE TABLE IF NOT EXISTS speakers (
    id SERIAL NOT NULL,
    name varchar(255) NOT NULL,
    bio text,
    event_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT speakers_pkey PRIMARY KEY (id),
    CONSTRAINT speakers_event_fk FOREIGN KEY (event_id) REFERENCES events (id)
    );

CREATE INDEX idx_speakers_deleted_at
    ON speakers USING btree
    (deleted_at ASC NULLS LAST);
