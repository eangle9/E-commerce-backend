CREATE TABLE inventory (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),      
product_item_id UUID NOT NULL,    
sku VARCHAR UNIQUE NOT NULL,                   
available_quantity INT DEFAULT 0,                
reserved_quantity INT DEFAULT 0,                  
minimum_quantity INT DEFAULT 0,     
status inventory_status NOT NULL DEFAULT 'IN_STOCK',              
restock_date TIMESTAMPTZ DEFAULT NULL,            
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),              
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),             
deleted_at TIMESTAMPTZ DEFAULT NULL              
);
CREATE INDEX idx_inventory_product_item_id ON inventory(product_item_id);
CREATE INDEX idx_inventory_sku ON inventory(sku);

ALTER TABLE inventory
ADD CONSTRAINT inventory_product_item_id_fkey
FOREIGN KEY (product_item_id) REFERENCES product_item(id) ON DELETE CASCADE;
