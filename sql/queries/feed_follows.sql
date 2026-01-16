-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT FF.*, U.name AS user_name, F.name AS feed_name
FROM inserted_feed_follow AS FF
INNER JOIN users AS U ON U.id = FF.user_id
INNER JOIN feeds AS F ON F.id = FF.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT FF.*, U.name AS user_name, F.name AS feed_name FROM feed_follows AS FF
INNER JOIN users AS U ON U.id = FF.user_id
INNER JOIN feeds AS F ON F.id = FF.feed_id
WHERE U.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;