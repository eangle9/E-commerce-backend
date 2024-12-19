CREATE TABLE wishlists (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL,
product_item_id UUID NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);

ALTER TABLE wishlists
ADD CONSTRAINT wishlists_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE wishlists
ADD CONSTRAINT wishlists_product_item_id_fkey
FOREIGN KEY (product_item_id) REFERENCES product_item(id) ON DELETE SET NULL;