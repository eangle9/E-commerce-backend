CREATE TABLE order_details (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL,
sub_total DECIMAL(10, 2) NOT NULL,
delivery_fee DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
service_charge DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
vat DECIMAL(10, 2) NOT NULL DEFAULT 0.00, 
total DECIMAL(10, 2) NOT NULL,
status order_status NOT NULL DEFAULT 'PENDING',
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_order_details_user_id ON order_details(user_id);

ALTER TABLE order_details
ADD CONSTRAINT order_details_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;