CREATE TABLE minecraft_server (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL,
    name TEXT NOT NULL,
    status INTEGER NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);

CREATE TABLE minecraft_server_port (
    server_id TEXT NOT NULL,
    port INTEGER NOT NULL,
    PRIMARY KEY (server_id, port)
);

-- TODO: add another table to handle members
