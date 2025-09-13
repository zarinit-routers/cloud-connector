-- +migrate Up
CREATE TABLE
    IF NOT EXISTS nodes (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        organization_id UUID NOT NULL,
        name VARCHAR(512),
        first_connection TIMESTAMPTZ NOT NULL,
        last_connection TIMESTAMPTZ
    );

CREATE TABLE
    IF NOT EXISTS tags (
        node_id UUID REFERENCES nodes (id) ON DELETE CASCADE NOT NULL,
        tag VARCHAR(64) NOT NULL,
        PRIMARY KEY (node_id, tag)
    );

-- +migrate Down
DROP TABLE tags;

DROP TABLE nodes;
