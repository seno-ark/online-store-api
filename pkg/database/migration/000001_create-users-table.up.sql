CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id"          uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "email"       VARCHAR(255) NOT NULL,
  "password"    VARCHAR(100) NOT NULL,
  "full_name"   VARCHAR(255) NOT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);