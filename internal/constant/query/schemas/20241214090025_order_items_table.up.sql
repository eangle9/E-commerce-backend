CREATE TABLE order_items (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
order_id UUID NOT NULL,
product_item_id UUID NOT NULL,
quantity INT NOT NULL,
total_price DECIMAL(10, 2) NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);

ALTER TABLE order_items
ADD CONSTRAINT order_items_order_id_fkey
FOREIGN KEY (order_id) REFERENCES order_details(id) ON DELETE CASCADE;
ALTER TABLE order_items
ADD CONSTRAINT order_items_product_item_id_fkey
FOREIGN KEY (product_item_id) REFERENCES product_item(id) ON DELETE CASCADE; 