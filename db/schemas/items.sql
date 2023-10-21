CREATE TYPE kind AS ENUM ('expenses', 'income');
CREATE TABLE items
(
    id          SERIAL PRIMARY KEY,
    user_id     UUID        NOT NULL,
    amount      INT         NOT NULL,
    kind        kind        NOT NULL,
    happened_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
