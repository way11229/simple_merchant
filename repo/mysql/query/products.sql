-- name: CreateProduct :execresult
INSERT INTO products (
    name,
    description,
    price,
    order_by,
    is_recommendation,
    total_quantity,
    sold_quantity,
    status
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ? 
)
;

-- name: ListTheRecommendedProducts :many
SELECT
    id,
    name,
    price,
    order_by
FROM
    products
WHERE
    1 = 1
    AND is_recommendation = true
    AND total_quantity >= sold_quantity
    AND status = 'on'
ORDER BY
    order_by DESC,
	created_at DESC
LIMIT ?
OFFSET ?
;

-- name: GetProductById :one
SELECT
    *
FROM
    products
WHERE
    1 = 1
    AND id = ? 
;

-- name: DeleteProductById :exec
DELETE FROM
    products
WHERE
    id = ? 
;