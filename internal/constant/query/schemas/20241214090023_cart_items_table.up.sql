CREATE TABLE cart_items (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
cart_id UUID NOT NULL,
product_item_id UUID NOT NULL,
price DECIMAL(10, 2) NOT NULL,
quantity INT NOT NULL,
total_price DECIMAL(10, 2) NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL    
);

ALTER TABLE cart_items
ADD CONSTRAINT cart_items_cart_id_fkey 
FOREIGN KEY (cart_id) REFERENCES shopping_cart(id);
ALTER TABLE cart_items
ADD CONSTRAINT cart_items_product_item_id_fkey
FOREIGN KEY (product_item_id) REFERENCES product_item(id);