-- name: ListComments :many
select * from list_comments_for_user(sqlc.arg(form_id), sqlc.arg(user_id));

-- name: CreateComment :one
select * from create_comment(
    sqlc.arg(form_id), sqlc.arg(user_id),
    sqlc.arg(body), sqlc.arg(element), sqlc.narg(parent)
);

-- name: UpdateComment :one
select * from update_comment(
    sqlc.arg(id), sqlc.arg(form_id), sqlc.arg(user_id),
    sqlc.narg(body), sqlc.narg(state)
);

-- name: DeleteComment :exec
select * from delete_comment(
    sqlc.arg(id), sqlc.arg(form_id), sqlc.arg(user_id)
);
