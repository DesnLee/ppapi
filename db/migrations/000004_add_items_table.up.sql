BEGIN;

--- 创建 kind 类型
CREATE TYPE kind AS ENUM ('expenses', 'income');
--- 创建 items 表
CREATE TABLE  IF NOT EXISTS items
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID        NOT NULL,
    amount      BIGINT      NOT NULL,
    kind        kind        NOT NULL,
    happened_at TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN items.user_id IS 'item belongs to user';
COMMENT ON COLUMN items.amount IS 'item amount';
COMMENT ON COLUMN items.kind IS 'item kind';
COMMENT ON COLUMN items.happened_at IS 'item happened time';


--- 创建 tags 表
CREATE TABLE  IF NOT EXISTS tags
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL,
    name       VARCHAR(50) NOT NULL,
    sign       VARCHAR(50) NOT NULL,
    kind       kind        NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN tags.user_id IS 'tag belongs to user';
COMMENT ON COLUMN tags.name IS 'tag name';
COMMENT ON COLUMN tags.sign IS 'tag sign';
COMMENT ON COLUMN tags.kind IS 'tag kind';
COMMENT ON COLUMN tags.deleted_at IS 'tag deleted time if not null';

--- 创建 items_tags 表
CREATE TABLE  IF NOT EXISTS items_tags
(
    id         BIGSERIAL PRIMARY KEY,
    item_id    BIGINT      NOT NULL,
    tag_id     BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN items_tags.item_id IS 'item id';
COMMENT ON COLUMN items_tags.tag_id IS 'tag id';

COMMIT;
