CREATE TABLE review (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
product_id UUID NOT NULL, 
user_id UUID NOT NULL,
order_id UUID NULL, 
rating rating NOT NULL, 
comment STRING,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_review_product_id ON review(product_id);

ALTER TABLE review
ADD CONSTRAINT review_product_id_fkey
FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE;
ALTER TABLE review
ADD CONSTRAINT review_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE review
ADD CONSTRAINT review_order_id_fkey
FOREIGN KEY (order_id) REFERENCES order_details(id) ON DELETE CASCADE;