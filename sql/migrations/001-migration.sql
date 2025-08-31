CREATE TABLE
    IF NOT EXISTS nodes (
        id UUID PRIMARY KEY,
        group_id UUID NOT NULL,
        title VARCHAR(512),
        first_connection TIMESTAMPTZ NOT NULL,
        last_connection TIMESTAMPTZ
    );

CREATE TABLE
    IF NOT EXISTS tags (
        node_id UUID REFERENCES nodes (id) ON DELETE CASCADE NOT NULL,
        title VARCHAR(64) NOT NULL,
        PRIMARY KEY (node_id, title)
    );
