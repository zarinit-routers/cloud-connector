-- name: GetNode :one
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

-- name: AddTag :exec

INSERT INTO tags (node_id, title) VALUES (@node_id,@title);

-- name: RemoveTag :exec

DELETE FROM tags t WHERE t.node_id = @node_id AND t.title = @title;

-- name: GetTags :many

SELECT t.title FROM tags t WHERE t.node_id = @node_id;
