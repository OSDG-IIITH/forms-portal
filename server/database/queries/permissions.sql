-- name: ListFormPermissions :many
select * from list_form_permissions(sqlc.arg(form_id), sqlc.arg(user_id));

-- name: GrantPermission :one
select * from grant_permission(
    sqlc.arg(form_id), sqlc.arg(user_id),
    sqlc.narg(target_user), sqlc.narg(target_group),
    sqlc.arg(role)::permission_role
);

-- name: RevokePermission :exec
select revoke_permission(
    sqlc.arg(form_id), sqlc.arg(user_id), sqlc.arg(permission_id)
);
