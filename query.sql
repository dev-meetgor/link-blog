-- Users

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = ?;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- Posts

-- name: CreatePost :one
INSERT INTO posts (author_id, url, title, content, slug)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPostByID :one
SELECT *
FROM posts
WHERE id = ?;

-- name: GetPostBySlug :one
SELECT *
FROM posts
WHERE slug = ?;

-- name: ListPostsByAuthor :many
SELECT *
FROM posts
WHERE author_id = ?
ORDER BY created_at DESC;

-- name: UpdatePost :exec
UPDATE posts
SET url = ?, title = ?, content = ?, slug = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = ?;

-- Tags

-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
RETURNING *;

-- name: GetTagByName :one
SELECT * FROM tags WHERE name = ?;

-- Post Tags

-- name: AddPostTag :exec
INSERT INTO post_tags (post_id, tag_id)
VALUES (?, ?);

-- name: RemovePostTag :exec
DELETE FROM post_tags
WHERE post_id = ? AND tag_id = ?;

-- name: GetPostTags :many
SELECT t.*
FROM tags t
JOIN post_tags pt ON t.id = pt.tag_id
WHERE pt.post_id = ?;
