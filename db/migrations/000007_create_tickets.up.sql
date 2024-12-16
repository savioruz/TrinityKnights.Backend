CREATE TABLE IF NOT EXISTS tickets (
    id varchar(36) NOT NULL,
    event_id integer NOT NULL,
    order_id integer DEFAULT NULL,
    price numeric(10,2) NOT NULL,
    type varchar(20) NOT NULL,
    seat_number varchar(10) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT tickets_pkey PRIMARY KEY (id),
    CONSTRAINT tickets_event_fk FOREIGN KEY (event_id) REFERENCES events (id),
    CONSTRAINT tickets_order_fk FOREIGN KEY (order_id) REFERENCES orders (id)
    );

ALTER TABLE tickets
    ADD CONSTRAINT tickets_seat_number_key UNIQUE (seat_number);

CREATE INDEX idx_tickets_deleted_at
    ON tickets USING btree
    (deleted_at ASC NULLS LAST);