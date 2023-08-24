-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: GetFeedsForUser :many
SELECT feeds.*
FROM feeds
where id IN (
        select feed_id
        from feed_follows
        where feed_follows.user_id = $1
    );
-- name: DeleteFeedFollow :exec 
DELETE FROM feed_follows
where id = $1
    AND user_id = $2;