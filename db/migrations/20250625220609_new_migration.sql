-- Create "accounts" table
CREATE TABLE "public"."accounts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "owner" text NOT NULL,
  "currency" text NOT NULL,
  "balance" numeric NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_accounts_deleted_at" to table: "accounts"
CREATE INDEX "idx_accounts_deleted_at" ON "public"."accounts" ("deleted_at");
-- Create "transfers" table
CREATE TABLE "public"."transfers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" numeric NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("from_account_id") REFERENCES "public"."accounts" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY ("to_account_id") REFERENCES "public"."accounts" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);
-- Create index "idx_transfers_deleted_at" to table: "transfers"
CREATE INDEX "idx_transfers_deleted_at" ON "public"."transfers" ("deleted_at");
-- Create "entries" table
CREATE TABLE "public"."entries" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "accounts_id" bigint NOT NULL,
  "amount" numeric NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("accounts_id") REFERENCES "public"."accounts" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);
-- Create index "idx_entries_deleted_at" to table: "entries"
CREATE INDEX "idx_entries_deleted_at" ON "public"."entries" ("deleted_at");
