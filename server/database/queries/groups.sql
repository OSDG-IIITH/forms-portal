-- name: ListGroups :many
select * from list_groups_for_user(
    sqlc.arg(user_id),
    sqlc.narg(owner_id),
    sqlc.narg(filter_type)::group_type,
    sqlc.arg(sort_by),
    sqlc.arg(order_by),
    sqlc.arg(limit_val),
    sqlc.arg(offset_val)
);

-- name: CountGroups :one
select count_groups_for_user(sqlc.arg(user_id));

-- name: CreateGroup :one
select * from create_group_of_type(
    sqlc.arg(owner_id),
    sqlc.arg(name),
    sqlc.narg(description),
    sqlc.arg(type),
    sqlc.narg(domain),
    sqlc.narg(members)::text[]
);

-- name: GetGroup :one
select * from get_group_by_id(sqlc.arg(id), sqlc.arg(user_id));

-- name: UpdateGroup :one
select * from update_group_by_id(
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.narg(name),
    sqlc.narg(description)
);

-- name: DeleteGroup :exec
select delete_group_by_id(sqlc.arg(id), sqlc.arg(user_id));

-- name: UpdateGroupDomain :exec
select update_domain_for_group(
    sqlc.arg(id), sqlc.arg(user_id), sqlc.arg(domain)
);

-- name: AddGroupMember :exec
select add_group_member_by_email(
    sqlc.arg(group_id), sqlc.arg(user_id), sqlc.arg(target_user)
);

-- name: RemoveGroupMember :exec
select remove_group_member_by_id(
    sqlc.arg(group_id), sqlc.arg(user_id), sqlc.arg(target_user_id)
);
