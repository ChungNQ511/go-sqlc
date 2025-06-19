-- name: CreateExample :one
INSERT INTO examples (example_name, example_description)
VALUES ($1, $2)
RETURNING *;

-- name: GetExample :one
SELECT * FROM examples
WHERE example_id = $1;

-- name: UpdateExample :one
UPDATE examples
SET example_name = $2, example_description = $3
WHERE example_id = $1 RETURNING *;

-- name: DeleteExample :exec
DELETE FROM examples
WHERE example_id = $1;
