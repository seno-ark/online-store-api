CREATE TABLE "payments" (
  "id"          		  uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "order_id"			    uuid UNIQUE NOT NULL,
  "payment_method" 		VARCHAR(100) NOT NULL,
  "payment_provider" 	VARCHAR(100) NOT NULL,
  "bill_amount" 		  INTEGER NOT NULL,
  "paid_amount" 		  INTEGER NOT NULL,
  "status"				    VARCHAR(20) NOT NULL,
  "transaction_id"	  VARCHAR(255) NOT NULL,
  "paid_at"  			    TIMESTAMPTZ NULL,
  "log"					      TEXT NULL,
  "created_at"  		  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"  		  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
CREATE INDEX ON "payments" ("order_id");