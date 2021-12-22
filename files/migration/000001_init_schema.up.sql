CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar NOT NULL UNIQUE,
  "fullname" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "items" (
  "id" bigserial PRIMARY KEY,
  "sku" varchar NOT NULL UNIQUE,
  "name" varchar NOT NULL,
  "stock" int NOT NULL,
  "price" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "status" int NOT NULL,
  "total" int NOT NULL,
  "expired_date" timestamp DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "quantity" int NOT NULL,
  "price" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "order_items" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

CREATE INDEX ON "items" ("sku");