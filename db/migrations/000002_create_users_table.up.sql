BEGIN;

CREATE TABLE users
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    email      VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN users.id IS 'unique identifier';
COMMENT ON COLUMN users.email IS 'unique email address';
COMMENT ON COLUMN users.created_at IS 'when the user was created';
COMMENT ON COLUMN users.updated_at IS 'when the user was last updated';

COMMIT;
