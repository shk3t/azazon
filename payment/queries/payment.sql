-- name: CreateProcessedPayment :exec
INSERT INTO processed_payment (order_id)
VALUES ($1);
