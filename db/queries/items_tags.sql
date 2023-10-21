-- name: CreateItemTagRelations :copyfrom
INSERT INTO items_tags (item_id, tag_id)
VALUES ($1, $2);
