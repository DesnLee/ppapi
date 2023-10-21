ALTER TABLE users
    ADD COLUMN name VARCHAR(42) NOT NULL DEFAULT '';

COMMENT ON COLUMN users.name IS 'user display name';
