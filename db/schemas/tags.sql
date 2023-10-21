CREATE TABLE tags
(
    id         SERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL,
    name       VARCHAR(50) NOT NULL,
    sign       VARCHAR(50) NOT NULL,
    kind       kind        NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
