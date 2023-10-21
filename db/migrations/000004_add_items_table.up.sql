BEGIN;

--- 创建 kind 类型
CREATE TYPE kind AS ENUM ('expenses', 'income');
--- 创建 items 表
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

--- 创建 tags 表
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

--- 创建 items_tags 表
CREATE TABLE items_tags
(
    id         SERIAL PRIMARY KEY,
    item_id    INT NOT NULL,
    tag_id     INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
