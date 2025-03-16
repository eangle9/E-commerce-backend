CREATE TABLE users (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
username VARCHAR(100) NOT NULL,
email VARCHAR(250) NOT NULL,
phone_number VARCHAR(50) NOT NULL,
password VARCHAR(100) NOT NULL,
first_name VARCHAR(50) NOT NULL,
last_name VARCHAR(50) NOT NULL,
profile_picture STRING NULL,
email_verified BOOL NOT NULL DEFAULT false,
role user_role NOT NULL DEFAULT 'CUSTOMER', 
last_login TIMESTAMPTZ NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
deleted_at TIMESTAMPTZ NULL
);
CREATE UNIQUE INDEX uni_idx_users_email ON users (email) WHERE deleted_at IS NULL;