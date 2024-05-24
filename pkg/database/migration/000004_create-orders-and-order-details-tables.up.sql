CREATE TABLE "orders" (
  "id"          uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "user_id" 	  uuid NOT NULL,
  "status"		  VARCHAR(20) NOT NULL,
  "other_cost" 	INTEGER NOT NULL,
  "total_cost" 	INTEGER NOT NULL,
  "shipment_address" VARCHAR(255) NOT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE INDEX ON "orders" ("user_id");
CREATE INDEX ON "orders" ("status");

CREATE TABLE "order_items" (
  "id"          	uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "order_id" 		uuid NOT NULL,
  "product_id" 		uuid NOT NULL,
  "notes"    		VARCHAR(100) NOT NULL,
  "qty" 			INTEGER NOT NULL,
  "product_price"	INTEGER NOT NULL,
  "created_at"  	TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

CREATE INDEX ON "order_items" ("order_id");
CREATE INDEX ON "order_items" ("product_id");