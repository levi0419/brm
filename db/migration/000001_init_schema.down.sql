--init_schema.down.sql

-- Drop foreign key constraints
ALTER TABLE "discounts" DROP CONSTRAINT IF EXISTS "fk_discounts_user_id";
ALTER TABLE "user_sessions" DROP CONSTRAINT IF EXISTS "fk_user_sessions_user_id";
ALTER TABLE "user_device_lists" DROP CONSTRAINT IF EXISTS "fk_user_device_lists_user_id";

-- Drop tables
DROP TABLE IF EXISTS "user_device_lists";
DROP TABLE IF EXISTS "user_sessions";
DROP TABLE IF EXISTS "discounts";
DROP TABLE IF EXISTS "users";
