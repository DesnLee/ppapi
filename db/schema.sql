CREATE TABLE users
(
    id         uuid PRIMARY KEY      DEFAULT gen_random_uuid(),
    email      VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);
