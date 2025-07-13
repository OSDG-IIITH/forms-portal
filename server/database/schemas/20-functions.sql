create or replace function has_form_permission(
    p_user_id text,
    p_form_id text,
    p_required_role permission_role
) returns boolean as $$
declare
    form_owner_id text;
    has_permission boolean;
begin
    select owner into form_owner_id from forms where id = p_form_id;
    if form_owner_id = p_user_id then
        return true;
    end if;

    select exists (
        select 1 from form_permissions
        where form = p_form_id and "user" = p_user_id and role = p_required_role
    ) into has_permission;

    return has_permission;
end;
$$ language plpgsql;

create or replace function list_forms_for_user(
    p_user_id text,
    p_filter_role permission_role,
    p_sort_by text,
    p_order text,
    p_limit int,
    p_offset int
) returns setof forms as $$
begin
    return query select f.* from forms f
    inner join form_permissions fp on f.id = fp.form and fp.user = p_user_id
    where p_filter_role is null or fp.role = p_filter_role
    group by f.id order by
        case when p_sort_by = 'modified' and p_order = 'asc' then f.modified end asc,
        case when p_sort_by = 'modified' and p_order = 'desc' then f.modified end desc,
        case when p_sort_by = 'title' and p_order = 'asc' then f.title end asc,
        case when p_sort_by = 'title' and p_order = 'desc' then f.title end desc
    limit p_limit offset p_offset;
end;
$$ language plpgsql;

create or replace function count_forms_for_user(
    p_user_id text,
    p_filter_role permission_role default null
) returns bigint as $$
declare
    v_count bigint;
begin
    select count(distinct f.id) into v_count from forms f
    left join form_permissions fp on f.id = fp.form and fp.user = p_user_id
    where p_filter_role is null or fp.role = p_filter_role;

    return v_count;
end;
$$ language plpgsql;

create or replace function create_form_with_permissions(
    p_owner_id text, p_slug text, p_title text, p_description text,
    p_structure text, p_live boolean, p_opens timestamptz, p_closes timestamptz,
    p_max_responses int, p_individual_limit int, p_editable_responses boolean
) returns forms as $$
declare
    v_form forms;
    v_perms permission_role[] := array['view', 'respond', 'comment', 'analyze', 'edit', 'manage'];
begin
    insert into forms (
        owner, slug, title, description, structure,
        live, opens, closes, max_responses, individual_limit, editable_responses
    ) values (
        p_owner_id, p_slug, p_title, p_description,
        p_structure, p_live, p_opens, p_closes,
        p_max_responses, p_individual_limit, p_editable_responses
    ) returning * into v_form;

    insert into form_permissions (form, "user", role)
    values (v_form.id, p_owner_id, unnest(v_perms));

    return v_form;
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

    if not has_form_permission(p_user_id, v_form.id, 'respond'::permission_role) then
        raise exception 'Form not found or you do not have permission to access it.' using hint = 'forbidden';
    end if;

    return v_form;
end;
$$ language plpgsql;

create or replace function get_form_by_id(
    p_id text,
    p_user_id text
) returns forms as $$
declare
    v_form forms;
begin
    select * into v_form from forms where id = p_id;

    if not found then
        raise exception 'Form not found or you do not have permission to access it.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, v_form.id, 'view'::permission_role) then
        raise exception 'Form not found or you do not have permission to access it.' using hint = 'forbidden';
    end if;

    return v_form;
end;
$$ language plpgsql;
