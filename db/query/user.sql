-- name: CreateUser :one
INSERT INTO users (
    uuid,
    first_name,
    last_name,
    phone_number,
    email,
    password,
    device_id,
    is_email_verified,
    created_at,
    status,
    last_updated_at,
    is_deleted
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, false, now(), false, DEFAULT, false
) RETURNING *;

-- name: GetUserByUUID :one
SELECT uuid, first_name, last_name, phone_number, email, created_at, last_updated_at
FROM users
WHERE uuid = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, uuid, first_name, last_name, phone_number, password, email, is_email_verified, created_at, last_updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUserPhoneNumber :one
UPDATE users 
SET phone_number = $2,
last_updated_at = now()
WHERE uuid = $1
RETURNING uuid, first_name, last_name, phone_number, email, created_at, last_updated_at;

-- name: UpdateUserEmail :one
UPDATE users 
SET email = $2,
last_updated_at = now()
WHERE uuid = $1
RETURNING uuid, first_name, last_name, phone_number, email, created_at, last_updated_at;

-- name: UpdateUserFirstName :one
UPDATE users 
SET first_name = $2,
last_updated_at = now()
WHERE uuid = $1
RETURNING uuid, first_name, last_name, phone_number, email, created_at, last_updated_at;

-- name: UpdateUserLastName :one
UPDATE users 
SET last_name = $2,
last_updated_at = now()
WHERE uuid = $1
RETURNING uuid, first_name, last_name, phone_number, email, created_at, last_updated_at;

-- name: UpdateUserDevice :one
UPDATE users 
SET device_id = $2,
last_updated_at = now()
WHERE uuid = $1
RETURNING uuid, first_name, last_name, phone_number, email, created_at, last_updated_at;

-- name: UpdateUserLastLogin :one
UPDATE users
SET  last_login = $2,
last_updated_at  = now()
WHERE uuid = $1
RETURNING uuid, first_name, last_name, phone_number, email, created_at, last_updated_at;
