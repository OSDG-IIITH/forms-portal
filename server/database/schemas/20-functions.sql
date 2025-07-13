create or replace function has_form_permission(
    p_user_id text,
    p_form_id text,
    p_required_roles permission_role[]
) returns boolean as $$
declare
    form_owner_id text;
    has_permission boolean;
begin
    select owner into form_owner_id from forms where id = p_form_id;
    if form_owner_id = p_user_id then return true; end if;

    select exists (
        select 1 from form_permissions
        where form = p_form_id
            and "user" = p_user_id
            and role = any(p_required_roles)
    ) into has_permission;

    return has_permission;
end;
$$ language plpgsql;

create or replace function resolve_form_by_handle_and_slug(
    p_handle text,
    p_slug text,
    p_user_id text
) returns forms as $$
declare
    v_form forms;
begin
    select f.* into v_form from forms f
    join users u on f.owner = u.id
    where u.handle = p_handle and f.slug = p_slug;

    if not found then
        raise exception 'Form not found or you do not have permission to access it.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, v_form.id, array['respond']::permission_role[]) then
        raise exception 'Form not found or you do not have permission to access it.' using hint = 'forbidden';
    end if;

    return v_form;
end;
$$ language plpgsql;
