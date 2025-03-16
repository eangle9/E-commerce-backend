CREATE TABLE product_items (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
product_id UUID NOT NULL,
color_id UUID DEFAULT NULL,
size_id UUID DEFAULT NULL,
sku VARCHAR(50) NOT NULL,
status product_status NOT NULL DEFAULT 'ACTIVE',
image_url STRING NOT NULL,
price DECIMAL(10, 2) NOT NULL,
discount DECIMAL(10, 2) NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMP NULL
);
CREATE UNIQUE INDEX uni_idx_product_items_sku ON product_items(sku);
CREATE INDEX idx_product_items_product_id ON product_items(product_id);

ALTER TABLE product_items
ADD CONSTRAINT product_items_product_id_fkey
FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE product_items
ADD CONSTRAINT product_items_color_id_fkey
FOREIGN KEY (color_id) REFERENCES color(id) ON DELETE SET NULL;
ALTER TABLE product_items
ADD CONSTRAINT product_items_size_id_fkey
FOREIGN KEY (size_id) REFERENCES size(id) ON DELETE SET NULL;


