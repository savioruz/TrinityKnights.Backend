CREATE TABLE IF NOT EXISTS payments (
    id SERIAL NOT NULL,
    order_id integer NOT NULL UNIQUE,
    method varchar(20),
    transaction_id varchar(255),
    amount float,
    status varchar(20),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT payments_pkey PRIMARY KEY (id),
    CONSTRAINT payments_order_fk FOREIGN KEY (order_id) REFERENCES orders (id)
    );

CREATE INDEX idx_payments_deleted_at
    ON payments USING btree
    (deleted_at ASC NULLS LAST);

ALTER TABLE orders
    ADD CONSTRAINT orders_payment_fk FOREIGN KEY (payment_id) REFERENCES payments (id);
