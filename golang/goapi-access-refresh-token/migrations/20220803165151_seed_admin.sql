-- +goose Up
-- +goose StatementBegin
-- user admin password admin
INSERT INTO "users" ("email","password","role") VALUES ('admin@mail.com','$argon2id$v=19$m=65536,t=3,p=4$5cR/it/OJUcqyRuqm4NaTA$pb8BUDHVCfw22kO3IsaqA7LoGock7an37KPOD9a49rE','admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "users" WHERE "email" = 'admin@mail.com';
-- +goose StatementEnd
