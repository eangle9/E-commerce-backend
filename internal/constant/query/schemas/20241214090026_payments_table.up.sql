CREATE TABLE payments (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
order_id UUID NOT NULL,
user_id UUID NOT NULL, 
payment_method VARCHAR(250) NULL, 
status payment_status NOT NULL DEFAULT 'PENDING', 
transaction_id VARCHAR UNIQUE NULL, 
tx_ref VARCHAR UNIQUE NOT NULL,
reference VARCHAR UNIQUE NULL,
type VARCHAR NULL,
total_amount DECIMAL(10, 2) NOT NULL,
currency_code currency_code NOT NULL DEFAULT 'ETB', 
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);

ALTER TABLE payments
ADD CONSTRAINT payments_order_id_fkey
FOREIGN KEY (order_id) REFERENCES order_details(id) ON DELETE CASCADE;
ALTER TABLE payments
ADD CONSTRAINT payments_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
