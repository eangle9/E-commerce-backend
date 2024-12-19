CREATE TABLE user_address (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL,
country STRING NOT NULL,
city STRING NOT NULL,
sub_city STRING NOT NULL,
woreda STRING NULL,
kebele STRING NULL,
street_address STRING NOT NULL,
phone_number STRING NOT NULL,
is_primary BOOL NOT NULL DEFAULT false,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE user_address
ADD CONSTRAINT user_address_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;