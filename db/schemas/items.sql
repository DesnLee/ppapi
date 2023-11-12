CREATE TABLE items
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID        NOT NULL,
    amount      BIGINT      NOT NULL,
    kind        VARCHAR(50) NOT NULL,
    happened_at TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
