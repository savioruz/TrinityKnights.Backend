CREATE TABLE IF NOT EXISTS orders (
    id SERIAL NOT NULL,
    user_id varchar(255) NOT NULL,
    date timestamp NOT NULL,
    total_price numeric(10,2) NOT NULL,
    payment_id integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT orders_pkey PRIMARY KEY (id),
    CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE INDEX idx_orders_deleted_at
    ON orders USING btree
    (deleted_at ASC NULLS LAST);