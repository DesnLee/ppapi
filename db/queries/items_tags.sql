-- name: CreateItemTagRelations :copyfrom
INSERT INTO items_tags (item_id, tag_id)
VALUES ($1, $2);

-- name: FindTagIDsByItemID :many
SELECT tag_id FROM items_tags WHERE item_id = $1;

-- name: DeleteAllItemTagRelation :exec
DELETE FROM items_tags;
