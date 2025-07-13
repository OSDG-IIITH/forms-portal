-- name: ListForms :many
select * from list_forms_for_user(
    sqlc.arg(user_id),
    sqlc.narg(filter_role)::permission_role,
    sqlc.arg(sort_by),
    sqlc.arg(order_by),
    sqlc.arg(limit_val),
    sqlc.arg(offset_val)
);

-- name: CountForms :one
select count_forms_for_user(
    sqlc.arg(user_id),
    sqlc.narg(filter_role)::permission_role
);

-- name: CreateFormWithPermissions :one
select * from create_form_with_permissions(
    sqlc.arg(owner_id),
    sqlc.arg(slug),
    sqlc.arg(title),
    sqlc.narg(description),
    sqlc.arg(structure),
    sqlc.arg(live),
    sqlc.narg(opens),
    sqlc.narg(closes),
    sqlc.narg(max_responses),
    sqlc.arg(individual_limit),
    sqlc.arg(editable_responses)
);

-- name: ResolveFormByHandleAndSlug :one
select * from resolve_form_by_handle_and_slug(
    sqlc.arg(handle), sqlc.arg(slug), sqlc.arg(user_id)
);
