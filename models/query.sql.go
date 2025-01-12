// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package models

import (
	"context"
	"database/sql"
)

const addPostTag = `-- name: AddPostTag :exec

INSERT INTO post_tags (post_id, tag_id)
VALUES (?, ?)
`

type AddPostTagParams struct {
	PostID int64
	TagID  int64
}

// Post Tags
func (q *Queries) AddPostTag(ctx context.Context, arg AddPostTagParams) error {
	_, err := q.db.ExecContext(ctx, addPostTag, arg.PostID, arg.TagID)
	return err
}

const createPost = `-- name: CreatePost :one

INSERT INTO posts (author_id, url, title, content, slug)
VALUES (?, ?, ?, ?, ?)
RETURNING id, author_id, url, title, content, slug, created_at, updated_at
`

type CreatePostParams struct {
	AuthorID int64
	Url      string
	Title    sql.NullString
	Content  sql.NullString
	Slug     sql.NullString
}

// Posts
func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.AuthorID,
		arg.Url,
		arg.Title,
		arg.Content,
		arg.Slug,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Url,
		&i.Title,
		&i.Content,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createTag = `-- name: CreateTag :one

INSERT INTO tags (name)
VALUES (?)
RETURNING id, name
`

// Tags
func (q *Queries) CreateTag(ctx context.Context, name string) (Tag, error) {
	row := q.db.QueryRowContext(ctx, createTag, name)
	var i Tag
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const createUser = `-- name: CreateUser :one

INSERT INTO users (username, email, password_hash)
VALUES (?, ?, ?)
RETURNING id, username, email, password_hash, created_at, updated_at
`

type CreateUserParams struct {
	Username     string
	Email        string
	PasswordHash string
}

// Users
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE id = ?
`

func (q *Queries) DeletePost(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPostByID = `-- name: GetPostByID :one
SELECT id, author_id, url, title, content, slug, created_at, updated_at
FROM posts
WHERE id = ?
`

func (q *Queries) GetPostByID(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByID, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Url,
		&i.Title,
		&i.Content,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostBySlug = `-- name: GetPostBySlug :one
SELECT id, author_id, url, title, content, slug, created_at, updated_at
FROM posts
WHERE slug = ?
`

func (q *Queries) GetPostBySlug(ctx context.Context, slug sql.NullString) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostBySlug, slug)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Url,
		&i.Title,
		&i.Content,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostTags = `-- name: GetPostTags :many
SELECT t.id, t.name
FROM tags t
JOIN post_tags pt ON t.id = pt.tag_id
WHERE pt.post_id = ?
`

func (q *Queries) GetPostTags(ctx context.Context, postID int64) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, getPostTags, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTagByName = `-- name: GetTagByName :one
SELECT id, name FROM tags WHERE name = ?
`

func (q *Queries) GetTagByName(ctx context.Context, name string) (Tag, error) {
	row := q.db.QueryRowContext(ctx, getTagByName, name)
	var i Tag
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE username = ?
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPostsByAuthor = `-- name: ListPostsByAuthor :many
SELECT id, author_id, url, title, content, slug, created_at, updated_at
FROM posts
WHERE author_id = ?
ORDER BY created_at DESC
`

func (q *Queries) ListPostsByAuthor(ctx context.Context, authorID int64) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, listPostsByAuthor, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.AuthorID,
			&i.Url,
			&i.Title,
			&i.Content,
			&i.Slug,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removePostTag = `-- name: RemovePostTag :exec
DELETE FROM post_tags
WHERE post_id = ? AND tag_id = ?
`

type RemovePostTagParams struct {
	PostID int64
	TagID  int64
}

func (q *Queries) RemovePostTag(ctx context.Context, arg RemovePostTagParams) error {
	_, err := q.db.ExecContext(ctx, removePostTag, arg.PostID, arg.TagID)
	return err
}

const updatePost = `-- name: UpdatePost :exec
UPDATE posts
SET url = ?, title = ?, content = ?, slug = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`

type UpdatePostParams struct {
	Url     string
	Title   sql.NullString
	Content sql.NullString
	Slug    sql.NullString
	ID      int64
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) error {
	_, err := q.db.ExecContext(ctx, updatePost,
		arg.Url,
		arg.Title,
		arg.Content,
		arg.Slug,
		arg.ID,
	)
	return err
}
