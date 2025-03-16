CREATE TABLE inventory (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),      
product_item_id UUID NOT NULL,    
sku VARCHAR(50) NOT NULL,                   
available_quantity INT8 NOT NULL DEFAULT 0,                
reserved_quantity INT8 DEFAULT 0,                  
minimum_quantity INT8 NOT NULL DEFAULT 1,     
status inventory_status NOT NULL DEFAULT 'IN_STOCK',              
restock_date TIMESTAMPTZ NULL,            
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),              
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),             
deleted_at TIMESTAMPTZ NULL              
);
CREATE UNIQUE INDEX uni_idx_inventory_sku ON inventory(sku);
CREATE INDEX idx_inventory_product_item_id ON inventory(product_item_id);
CREATE INDEX idx_inventory_sku ON inventory(sku);

ALTER TABLE inventory
ADD CONSTRAINT inventory_product_item_id_fkey
FOREIGN KEY (product_item_id) REFERENCES product_items(id) ON DELETE CASCADE;
