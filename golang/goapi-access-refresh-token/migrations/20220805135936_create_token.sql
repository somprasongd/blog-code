-- +goose Up
-- +goose StatementBegin

CREATE TABLE tokens (
    id      uuid NOT NULL DEFAULT uuid_generate_v4(),
    tokenId uuid NOT NULL,
    userId  uuid NOT NULL,
    created_at timestamptz NOT NULL default current_timestamp,
    expires_at timestamptz NOT NULL,
	  CONSTRAINT tokens_pkey PRIMARY KEY (id),
    CONSTRAINT tokens_userId_fkey FOREIGN KEY (userId) REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE
);
CREATE UNIQUE INDEX tokens_token_key ON tokens(tokenId);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;
-- +goose StatementEnd
