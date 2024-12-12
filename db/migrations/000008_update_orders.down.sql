BEGIN;

ALTER TABLE orders 
    ADD COLUMN payment_id integer;

ALTER TABLE orders
    ADD CONSTRAINT orders_payment_fk FOREIGN KEY (payment_id) REFERENCES payments (id);

COMMIT;
