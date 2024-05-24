CREATE TABLE "categories" (
  "id"          uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "name"    	VARCHAR(50) NOT NULL,
  "description" VARCHAR(255) NOT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "products" (
  "id"          uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "category_id" uuid NOT NULL,
  "name"    	VARCHAR(100) NOT NULL,
  "description" TEXT NOT NULL,
  "price" 		INTEGER NOT NULL,
  "stock" 		INTEGER NOT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");
CREATE INDEX ON "products" ("category_id");