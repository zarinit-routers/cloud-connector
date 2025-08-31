-- name: NodeExists :one
SELECT
    *
FROM
    nodes n
WHERE
    n.id = @id
LIMIT
    1;

-- name: NewNode :exec
INSERT INTO nodes (id, group_id, first_connection, last_connection) 
VALUES (@id, @group_id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: GetNodes :many
SELECT
    *
FROM
    nodes n
WHERE
    n.group_id = @group_id;

-- name: UpdateLastConnection :exec
UPDATE
    nodes n
SET
    last_connection = CURRENT_TIMESTAMP
WHERE
    n.id = @id;


-- name: ReconnectNode :exec
UPDATE
    nodes n
SET
    last_connection = CURRENT_TIMESTAMP,
    group_id = @group_id
WHERE
    n.id = @id;
