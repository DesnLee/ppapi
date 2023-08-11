CREATE TABLE users
(
    id         uuid PRIMARY KEY             DEFAULT gen_random_uuid(),
    email      VARCHAR(255) NOT NULL UNIQUE,
    phone      VARCHAR(50)  NOT NULL UNIQUE,
    address    VARCHAR(255),
    created_at TIMESTAMPTZ  NOT NULL        DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL        DEFAULT now()
);
