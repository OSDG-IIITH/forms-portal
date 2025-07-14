-- name: ListGroups :many
select * from list_groups_for_user(
    sqlc.arg(user_id),
    sqlc.arg(sort_by),
    sqlc.arg(order_by),
    sqlc.arg(limit_val),
    sqlc.arg(offset_val)
);

-- name: CountGroups :one
select count_groups_for_user(sqlc.arg(user_id));

-- name: CreateGroup :one
select * from create_group(
    sqlc.arg(owner_id),
    sqlc.arg(name),
    sqlc.narg(description),
    sqlc.arg(type),
    sqlc.narg(domain),
    sqlc.narg(members)::text[]
);
