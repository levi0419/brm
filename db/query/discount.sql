-- name: AddDiscount :one
INSERT INTO discounts (
    "user_id",
    "label",
    "expiration_time",
    "code",
    "status",
    "created_at"
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;



-- name: GetUserDiscountByUserID :one
SELECT * FROM discounts
WHERE user_id = $1 LIMIT 1;
