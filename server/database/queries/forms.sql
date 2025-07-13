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
