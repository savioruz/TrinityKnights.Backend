BEGIN;

ALTER TABLE orders 
    DROP CONSTRAINT IF EXISTS orders_payment_fk;

ALTER TABLE orders 
    DROP COLUMN IF EXISTS payment_id;

COMMIT;
