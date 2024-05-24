CREATE TABLE "carts" (
  "id"          uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "user_id" 	uuid NOT NULL,
  "product_id" 	uuid NOT NULL,
  "notes"    	VARCHAR(100) NOT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "carts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "carts" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
ALTER TABLE "carts" ADD CONSTRAINT "user_product_key" UNIQUE ("user_id", "product_id");

CREATE INDEX ON "carts" ("user_id");
CREATE INDEX ON "carts" ("product_id");
