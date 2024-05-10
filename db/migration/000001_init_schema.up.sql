CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY NOT NULL,
  "uuid" uuid NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "device_id" varchar NOT NULL,
  "last_login" timestamp,
  "login_attempt" int,
  "status" varchar,
  "is_email_verified"  VARCHAR NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "last_updated_at" timestamp DEFAULT (now()),
  "is_deleted" boolean DEFAULT false,
  "deleted_at" timestamp,
  "deleted_by" bigint
);

CREATE TABLE "discounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "label" varchar NOT NULL,
  "expiration_time" timestamp NOT NULL,
  "code" VARCHAR NOT NULL,
  "status" boolean DEFAULT false NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "user_sessions" (
  "id" BIGSERIAL PRIMARY KEY,
  "session_id" uuid NOT NULL,
  "user_id" bigint  NOT NULL,
  "token" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "ip" varchar NOT NULL,
  "channel" varchar NOT NULL,
  -- "country" varchar NOT NULL,
  -- "state" varchar NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "user_device_lists" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "device_id" varchar NOT NULL,
  "ip_address" varchar NOT NULL,
  "client_details" varchar NOT NULL,
  "is_blocked" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "discounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_device_lists" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
