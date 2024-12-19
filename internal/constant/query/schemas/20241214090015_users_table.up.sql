CREATE TABLE users (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
username STRING NOT NULL,
email STRING NOT NULL,
phone_number STRING NOT NULL,
password STRING NOT NULL,
first_name STRING NOT NULL,
last_name STRING,
profile_picture STRING NULL,
email_verified BOOL NOT NULL DEFAULT false,
role user_role NOT NULL DEFAULT 'CUSTOMER', 
last_login TIMESTAMPTZ DEFAULT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
deleted_at TIMESTAMPTZ NULL
);
CREATE UNIQUE INDEX idx_users_email on users (email);