CREATE TABLE "todos" (
  "id" BIGSERIAL NOT NULL,
  "title" VARCHAR NOT NULL,
  "is_done" BOOLEAN NOT NULL DEFAULT false,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
  PRIMARY KEY ("id")
);