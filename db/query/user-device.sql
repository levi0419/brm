-- name: AddUserDevice :one
INSERT INTO user_device_lists (
    user_id,
    device_id,
    ip_address,
    client_details
) VALUES (
    $1, $2, $3, $4
)RETURNING *;