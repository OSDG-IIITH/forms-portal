-- name: ListFormPermissions :many
select * from list_permissions_for_form(sqlc.arg(form_id), sqlc.arg(user_id));

-- name: GrantPermission :one
select * from grant_permission_on_form(
    sqlc.arg(form_id), sqlc.arg(user_id),
    sqlc.narg(target_user), sqlc.narg(target_group),
    sqlc.arg(role)::permission_role
);

-- name: RevokePermission :exec
select revoke_permission_by_id(
    sqlc.arg(form_id), sqlc.arg(user_id), sqlc.arg(permission_id)
);
