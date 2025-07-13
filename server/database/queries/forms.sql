-- name: ListForms :many
select forms.* from forms
left join form_permissions on forms.id = form_permissions.form
where
  forms.owner = sqlc.arg(owner_id)
  or (
    form_permissions.user = sqlc.arg(owner_id)
    and (
      sqlc.narg(role)::permission_role is null
      or form_permissions.role = sqlc.narg(role)::permission_role
    )
  )
limit sqlc.arg(limit_val) offset sqlc.arg(offset_val);

-- name: CountForms :one
select count(distinct forms.id)
from forms
left join form_permissions on forms.id = form_permissions.form
where
  forms.owner = sqlc.arg(owner_id)
  or (
    form_permissions.user = sqlc.arg(owner_id)
    and (
      sqlc.narg(role)::permission_role is null
      or form_permissions.role = sqlc.narg(role)::permission_role
    )
  )
;

-- name: CreateForm :one
insert into forms (
    owner, slug, title, description, structure,
    live, opens, closes, max_responses, individual_limit, editable_responses
) values (
    sqlc.arg(owner_id), sqlc.arg(slug), sqlc.arg(title), sqlc.narg(description),
    sqlc.arg(structure), sqlc.arg(live), sqlc.narg(opens), sqlc.narg(closes),
    sqlc.narg(max_responses), sqlc.arg(individual_limit), sqlc.arg(editable_responses)
) returning *;

-- name: GetFormByHandleAndSlug :one
select forms.* from forms
join users on forms.owner = users.id
where users.handle = sqlc.arg(handle) and forms.slug = sqlc.arg(slug);
