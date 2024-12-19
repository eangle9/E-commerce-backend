CREATE TABLE shopping_cart (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL, 
sub_total DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
delivery_fee DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
service_charge DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
vat DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
total DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL   
);

ALTER TABLE shopping_cart
ADD CONSTRAINT shopping_cart_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id);