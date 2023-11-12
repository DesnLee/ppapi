BEGIN;

CREATE TABLE IF NOT EXISTS validation_codes
(
    id         BIGSERIAL PRIMARY KEY,
    code       VARCHAR(20)  NOT NULL,
    email      VARCHAR(255) NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN validation_codes.code IS 'email validation code';
COMMENT ON COLUMN validation_codes.email IS 'email address';
COMMENT ON COLUMN validation_codes.used_at IS 'when the code was used';
COMMENT ON COLUMN validation_codes.created_at IS 'when the code was created';
COMMENT ON COLUMN validation_codes.updated_at IS 'when the code was last updated';

COMMIT;
