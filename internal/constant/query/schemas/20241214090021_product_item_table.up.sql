CREATE TABLE product_item (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
product_id UUID NOT NULL,
color_id UUID DEFAULT NULL,
size_id UUID DEFAULT NULL,
-- reserved_quantity INT NOT NULL DEFAULT 0,
sku STRING UNIQUE NOT NULL,
status product_status NOT NULL DEFAULT 'ACTIVE',
image_url STRING NOT NULL,
price DECIMAL(10, 2) NOT NULL,
discount DECIMAL(10, 2) DEFAULT NULL,
-- qty_in_stock INT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_product_item_product_id ON product_item(product_id);

ALTER TABLE product_item
ADD CONSTRAINT product_item_product_id_fkey
FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE;
ALTER TABLE product_item
ADD CONSTRAINT product_item_color_id_fkey
FOREIGN KEY (color_id) REFERENCES color(id) ON DELETE SET NULL;
ALTER TABLE product_item
ADD CONSTRAINT product_item_size_id_fkey
FOREIGN KEY (size_id) REFERENCES size(id) ON DELETE SET NULL;


