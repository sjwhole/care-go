-- name: GetSubscriptionById :one
SELECT *
FROM subscriptions
WHERE id = ?;

-- name: GetSubscriptionsByUserId :many
SELECT *
FROM subscriptions
WHERE user_id = ?
ORDER BY expires_at DESC;

-- name: CreateSubscription :execresult
INSERT INTO subscriptions (user_id, expires_at)
VALUES (?, ?);


