-- +goose Up
-- +goose StatementBegin
-- user admin password admin
INSERT INTO "users" ("email","password","role") VALUES ('admin@mail.com','$2a$14$rgdCRN582Nu5rKgfQA/TiOVb.YIx.SzZjE6NZW8AErhcIIajsnZLy','admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "users" WHERE "email" = 'admin@mail.com';
-- +goose StatementEnd
