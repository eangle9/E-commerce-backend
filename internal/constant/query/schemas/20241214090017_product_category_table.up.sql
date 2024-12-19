CREATE TABLE product_category (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
parent_id UUID DEFAULT NULL,
name STRING NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE product_category
ADD CONSTRAINT product_category_parent_id_fkey
FOREIGN KEY (parent_id) REFERENCES product_category(id) ON DELETE SET NULL;