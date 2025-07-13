-- name: ListFormPermissions :many
select * from list_form_permissions(
    sqlc.arg(form_id), sqlc.arg(user_id)
);
