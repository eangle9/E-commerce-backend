CREATE TABLE user_address (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL,
country VARCHAR(250) NOT NULL,
city VARCHAR(250) NOT NULL,
sub_city VARCHAR(250) NOT NULL,
woreda VARCHAR(250) NULL,
kebele VARCHAR(50) NULL,
street VARCHAR(250) NOT NULL,
phone_number VARCHAR(50) NOT NULL,
is_primary BOOL NOT NULL DEFAULT false,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ NULL
);

ALTER TABLE user_address
ADD CONSTRAINT user_address_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;