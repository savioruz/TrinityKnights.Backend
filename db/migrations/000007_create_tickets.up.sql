CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL NOT NULL,
    event_id integer NOT NULL,
    order_id integer NOT NULL,
    price numeric(10,2) NOT NULL,
    type varchar(20) NOT NULL,
    seat_number integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT tickets_pkey PRIMARY KEY (id),
    CONSTRAINT tickets_event_fk FOREIGN KEY (event_id) REFERENCES events (id),
    CONSTRAINT tickets_order_fk FOREIGN KEY (order_id) REFERENCES orders (id)
    );

CREATE INDEX idx_tickets_deleted_at
    ON tickets USING btree
    (deleted_at ASC NULLS LAST);