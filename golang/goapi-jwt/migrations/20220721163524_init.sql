-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE public.users (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	email varchar(255) NOT NULL,
	password varchar(255) NOT NULL,
	role user_role NULL DEFAULT 'user'::user_role,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX users_unique_email ON public.users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TYPE user_role;
-- +goose StatementEnd
