CREATE TABLE product (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
category_id UUID NOT NULL,
brand VARCHAR NULL,
name VARCHAR NOT NULL,
status product_status NOT NULL DEFAULT 'ACTIVE',
description STRING NULL,
average_rating DECIMAL(3, 2) NOT NULL DEFAULT 0.00,
total_reviews INT NOT NULL DEFAULT 0,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_product_category_id ON product(category_id);

ALTER TABLE product
ADD CONSTRAINT product_category_parent_id_fkey
FOREIGN KEY (category_id) REFERENCES product_category(id) ON DELETE CASCADE;