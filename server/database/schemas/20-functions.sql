create or replace function has_form_permission(
    p_user_id text,
    p_form_id text,
    p_required_role permission_role
) returns boolean as $$
begin
    return exists (
        select 1 from form_permissions
        where
            form = p_form_id and
            role = p_required_role and
            "user" = p_user_id

        union all

        select 1 from form_permissions as fp
        join groups as g on fp."group" = g.id
        join group_list_members as glm on g.id = glm."group"
        where
            fp.form = p_form_id and
            fp.role = p_required_role and
            g.type = 'list' and
            glm."user" = p_user_id

        union all

        select 1 from form_permissions as fp
        join groups as g on fp."group" = g.id
        join group_domain_rules as gdr on g.id = gdr."group"
        where
            fp.form = p_form_id and
            fp.role = p_required_role and
            g.type = 'domain' and
            gdr.domain = (
                select substring(email from '@(.*)$')
                from users
                where id = p_user_id
            )
    );
end;
$$ language plpgsql;

create or replace function has_group_permission(
    p_user_id text,
    p_group_id text,
    p_required_type group_type
) returns boolean as $$
declare
    has_permission boolean;
begin
    select exists (
        select 1 from groups
        where id = p_group_id and owner = p_user_id and (
            type is null or type = p_required_type
        )
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
    p_anonymous boolean, p_max_responses int, p_individual_limit int,
    p_editable_responses boolean
) returns forms as $$
declare
    v_form forms;
    v_perms permission_role[] := array['view', 'respond', 'comment', 'analyze', 'edit', 'manage'];
begin
    insert into forms (
        owner, slug, title, description, structure,
        live, opens, closes, anonymous, max_responses, individual_limit,
        editable_responses
    ) values (
        p_owner_id, p_slug, p_title, p_description,
        p_structure, p_live, p_opens, p_closes, p_anonymous,
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
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, v_form.id, 'respond'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
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
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, v_form.id, 'view'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    return v_form;
end;
$$ language plpgsql;

create or replace function update_form_by_id(
    p_id text, p_user_id text, p_slug text, p_title text, p_description text,
    p_structure text, p_live boolean, p_opens timestamptz, p_closes timestamptz,
    p_anonymous boolean, p_max_responses int, p_individual_limit int,
    p_editable_responses boolean
) returns forms as $$
declare
    v_form forms;
begin
    if not has_form_permission(p_user_id, p_id, 'edit'::permission_role) then
        raise exception 'You do not have permission to edit this form.' using hint = 'forbidden';
    end if;

    update forms set
        slug = coalesce(p_slug, slug), title = coalesce(p_title, title),
        description = coalesce(p_description, description),
        structure = coalesce(p_structure, structure), live = coalesce(p_live, live),
        opens = coalesce(p_opens, opens), closes = coalesce(p_closes, closes),
        anonymous = coalesce(p_anonymous, anonymous),
        max_responses = coalesce(p_max_responses, max_responses),
        individual_limit = coalesce(p_individual_limit, individual_limit),
        editable_responses = coalesce(p_editable_responses, editable_responses)
    where id = p_id
    returning * into v_form;

    return v_form;
end;
$$ language plpgsql;

create or replace function delete_form_by_id(
    p_id text,
    p_user_id text
) returns void as $$
begin
    if not has_form_permission(p_user_id, p_id, 'manage'::permission_role) then
        raise exception 'You do not have permission to delete this form.' using hint = 'forbidden';
    end if;

    delete from forms where id = p_id;
end;
$$ language plpgsql;

create or replace function list_permissions_for_form(
    p_form_id text,
    p_user_id text
) returns setof form_permissions as $$
begin
    if not has_form_permission(p_user_id, p_form_id, 'manage'::permission_role) then
        raise exception 'You do not have permission to manage this form.' using hint = 'forbidden';
    end if;

    return query select * from form_permissions where form = p_form_id;
end;
$$ language plpgsql;

create or replace function grant_permission_on_form(
    p_form_id text,
    p_user_id text,
    p_target_user text,
    p_target_group text,
    p_role permission_role
) returns setof form_permissions as $$
declare
    v_target_user_id text;
begin
    if not has_form_permission(p_user_id, p_form_id, 'manage'::permission_role) then
        raise exception 'You do not have permission to manage this form.' using hint = 'forbidden';
    end if;

    select u.id into v_target_user_id from users u where u.email = p_target_user;
    if not found and p_target_user is not null then
        raise exception 'User with email % does not exist.', p_target_user using hint = 'not-found';
    end if;

    return query select * from form_permissions
    where form = p_form_id and role = p_role
      and "user" is not distinct from v_target_user_id
      and "group" is not distinct from p_target_group;

    if not found then
        return query insert into form_permissions (form, role, "user", "group")
        values (p_form_id, p_role, v_target_user_id, p_target_group) returning *;
    end if;
end;
$$ language plpgsql;

create or replace function revoke_permission_by_id(
    p_form_id text, p_user_id text, p_permission_id text
) returns void as $$
begin
    if not has_form_permission(p_user_id, p_form_id, 'manage'::permission_role) then
        raise exception 'You do not have permission to manage this form.' using hint = 'forbidden';
    end if;

    delete from form_permissions where id = p_permission_id;
end;
$$ language plpgsql;

create or replace function list_groups_for_user(
    p_user_id text,
    p_owner_email text,
    p_filter_type group_type,
    p_sort_by text,
    p_order text,
    p_limit int,
    p_offset int
) returns setof group_with_details as $$
begin
    return query with group_details as (
        select g.*, d.domain, array_agg(m."user" order by m."user")
        filter (where m."user" is not null) as members from groups g
        left join group_domain_rules d on g.id = d."group"
        left join group_list_members m on g.id = m."group"
        where (p_owner_email is null or owner = (
            select id from users where email = p_owner_email
        )) group by g.id, d.domain
    ) select * from group_details where owner = p_user_id or (
        p_user_id = any(members) or domain = (
            select substring(email from '@(.*)$') from users where id = p_user_id
        )
    ) and (p_filter_type is null or type = p_filter_type) order by
        case when p_sort_by = 'created' and p_order = 'asc' then id end asc,
        case when p_sort_by = 'created' and p_order = 'desc' then id end desc,
        case when p_sort_by = 'name' and p_order = 'asc' then name end asc,
        case when p_sort_by = 'name' and p_order = 'desc' then name end desc,
        case when p_sort_by = 'type' and p_order = 'asc' then type end asc,
        case when p_sort_by = 'type' and p_order = 'desc' then type end desc
    limit p_limit offset p_offset;
end;
$$ language plpgsql;

create or replace function count_groups_for_user(p_user_id text)
returns bigint as $$
declare
    v_count bigint;
begin
    select count(distinct id) into v_count from groups where owner = p_user_id;

    return v_count;
end;
$$ language plpgsql;

create or replace function create_group_of_type(
    p_owner_id text,
    p_name text,
    p_description text,
    p_type group_type,
    p_domain text,
    p_members text[]
) returns group_with_details as $$
declare
    v_group_id text;
    v_missing_emails text[];
    v_group group_with_details;
begin
    insert into groups (owner, name, description, type)
    values (p_owner_id, p_name, p_description, p_type)
    returning groups.id into v_group_id;

    if p_type = 'domain' and p_domain is not null then
        insert into group_domain_rules ("group", domain)
        values (v_group_id, p_domain);
    end if;

    if p_type = 'list' and p_members is not null then
        select array_agg(i_email) into v_missing_emails
        from unnest(p_members) as i_email
        left join users u on u.email = i_email where u.id is null;

        if v_missing_emails is not null then
            raise exception 'Could not create group with users that do not exist: %',
                array_to_string(v_missing_emails, ', ') using hint = 'not-found';
        end if;

        insert into group_list_members ("group", "user")
        select v_group_id, u.id from unnest(p_members) as i_email
        join users u on u.email = i_email on conflict do nothing;
    end if;

    select g.*, d.domain, array_agg(m."user" order by m."user")
    filter (where m."user" is not null) as members into v_group from groups g
    left join group_domain_rules d on g.id = d."group"
    left join group_list_members m on g.id = m."group"
    where g.id = v_group_id
    group by g.id, g.owner, g.name, g.description, g.type, d.domain;

    return v_group;
end;
$$ language plpgsql;

create or replace function get_group_by_id(
    p_id text,
    p_user_id text
) returns group_with_details as $$
declare
    v_group group_with_details;
begin
    if not has_group_permission(p_user_id, p_id) then
        raise exception 'Group not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    select g.*, d.domain, array_agg(m."user" order by m."user")
    filter (where m."user" is not null) as members into v_group from groups g
    left join group_domain_rules d on g.id = d."group"
    left join group_list_members m on g.id = m."group"
    where g.id = p_id and g.owner = p_user_id
    group by g.id, g.owner, g.name, g.description, g.type, d.domain;

    if not found then
        raise exception 'Group not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    return v_group;
end;
$$ language plpgsql;

create or replace function update_group_by_id(
    p_id text,
    p_user_id text,
    p_name text,
    p_description text
) returns group_with_details as $$
declare
    v_group group_with_details;
begin
    if not has_group_permission(p_user_id, p_id) then
        raise exception 'Group not found or you do not have permission to do this.' using hint = 'forbidden';
    end if;

    update groups set
        name = coalesce(p_name, name),
        description = coalesce(p_description, description)
    where id = p_id and owner = p_user_id
    returning * into v_group;

    return v_group;
end;
$$ language plpgsql;

create or replace function update_domain_for_group(
    p_id text,
    p_user_id text,
    p_domain text
) returns void as $$
begin
    if not has_group_permission(p_user_id, p_id, 'domain'::group_type) then
        raise exception 'Group not found or you do not have permission to do this.' using hint = 'forbidden';
    end if;

    update group_domain_rules set domain = p_domain where "group" = p_id;
end;
$$ language plpgsql;

create or replace function delete_group_by_id(
    p_id text,
    p_user_id text
) returns void as $$
begin
    if not has_group_permission(p_user_id, p_id) then
        raise exception 'Group not found or you do not have permission to do this.' using hint = 'forbidden';
    end if;

    delete from groups where id = p_id;
end;
$$ language plpgsql;

create or replace function add_group_member_by_email(
    p_group_id text,
    p_user_id text,
    p_target_user text
) returns void as $$
declare
    v_target_user_id text;
begin
    if not has_group_permission(p_user_id, p_group_id, 'list'::group_type) then
        raise exception 'Group not found or you do not have permission to do this.' using hint = 'forbidden';
    end if;

    select u.id into v_target_user_id from users u where u.email = p_target_user;
    if not found then
        raise exception 'User with email % does not exist.', p_target_user using hint = 'not-found';
    end if;

    insert into group_list_members ("group", "user")
    values (p_group_id, v_target_user_id) on conflict do nothing;
end;
$$ language plpgsql;

create or replace function remove_group_member_by_id(
    p_group_id text,
    p_user_id text,
    p_target_user_id text
) returns void as $$
begin
    if not has_group_permission(p_user_id, p_group_id, 'list'::group_type) then
        raise exception 'Group not found or you do not have permission to do this.' using hint = 'forbidden';
    end if;

    delete from group_list_members
    where "group" = p_group_id and "user" = p_target_user_id;
end;
$$ language plpgsql;

create or replace function list_comments_for_form(
    p_form_id text,
    p_user_id text
) returns setof comments as $$
begin
    if not has_form_permission(p_user_id, p_form_id, 'comment'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    return query select * from comments where form = p_form_id
    order by modified desc;
end;
$$ language plpgsql;

create or replace function create_comment_on_form(
    p_form_id text, p_user_id text,
    p_body text, p_element text, p_parent text
) returns comments as $$
declare
    v_comment comments;
begin
    if not has_form_permission(p_user_id, p_form_id, 'comment'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    insert into comments (form, commenter, body, element, parent)
    values (p_form_id, p_user_id, p_body, p_element, p_parent)
    returning * into v_comment;

    return v_comment;
end;
$$ language plpgsql;

create or replace function update_comment_by_id(
    p_id text, p_form_id text, p_user_id text,
    p_body text, p_state comment_state
) returns comments as $$
declare
    v_comment comments;
begin
    if not has_form_permission(p_user_id, p_form_id, 'comment'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    select * into v_comment from comments where id = p_id and form = p_form_id;
    if not found then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if v_comment.commenter != p_user_id then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    update comments set
        body = coalesce(p_body, body),
        state = coalesce(p_state, state)
    where id = p_id returning * into v_comment;

    return v_comment;
end;
$$ language plpgsql;

create or replace function delete_comment_by_id(
    p_id text, p_form_id text, p_user_id text
) returns void as $$
declare
    v_comment comments;
begin
    select * into v_comment from comments where id = p_id and form = p_form_id;

    if not found then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if v_comment.commenter != p_user_id and not has_form_permission(p_user_id, p_form_id, 'manage'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    delete from comments where id = p_id;
end;
$$ language plpgsql;

create or replace function list_responses_for_form(
    p_form_id text,
    p_user_id text,
    p_status response_status,
    p_sort_by text,
    p_order text,
    p_limit int,
    p_offset int
) returns setof responses as $$
begin
    if has_form_permission(p_user_id, p_form_id, 'analyze'::permission_role) then
        return query select * from responses r where r.form = p_form_id
        and r.status = p_status order by
            case when p_sort_by = 'submitted' and p_order = 'asc' then r.submitted end asc,
            case when p_sort_by = 'submitted' and p_order = 'desc' then r.submitted end desc,
            case when p_sort_by = 'started' and p_order = 'asc' then r.started end asc,
            case when p_sort_by = 'started' and p_order = 'desc' then r.started end desc,
            case when p_sort_by = 'edited' and p_order = 'asc' then r.edited end asc,
            case when p_sort_by = 'edited' and p_order = 'desc' then r.edited end desc
        limit p_limit offset p_offset;

        return;
    end if;

    if has_form_permission(p_user_id, p_form_id, 'respond'::permission_role) then
        return query select * from responses r where r.form = p_form_id
        and r.status = p_status and r.respondent = p_user_id order by
            case when p_sort_by = 'submitted' and p_order = 'asc' then r.submitted end asc,
            case when p_sort_by = 'submitted' and p_order = 'desc' then r.submitted end desc,
            case when p_sort_by = 'started' and p_order = 'asc' then r.started end asc,
            case when p_sort_by = 'started' and p_order = 'desc' then r.started end desc,
            case when p_sort_by = 'edited' and p_order = 'asc' then r.edited end asc,
            case when p_sort_by = 'edited' and p_order = 'desc' then r.edited end desc
        limit p_limit offset p_offset;

        return;
    end if;

    raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
end;
$$ language plpgsql;

create or replace function count_responses_for_form(
    p_form_id text,
    p_user_id text,
    p_status response_status
) returns bigint as $$
declare
    v_count bigint;
begin
    if not has_form_permission(p_user_id, p_form_id, 'analyze'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    select count(distinct r.id) into v_count from responses r
    where r.form = p_form_id and r.status = p_status;

    return v_count;
end;
$$ language plpgsql;

create or replace function start_response_for_form(
    p_form_id text,
    p_user_id text
) returns responses as $$
declare
    v_form forms;
    v_response responses;
    v_existing_count int;
    v_total_count int;
begin
    if not has_form_permission(p_user_id, p_form_id, 'respond'::permission_role) then
        raise exception 'Form not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    select * into v_form from forms f where f.id = p_form_id;

    select sr.responses into v_existing_count from submission_records sr
    where sr.form = p_form_id and sr."user" = p_user_id;
    select sum(sr.responses) into v_total_count from submission_records sr
    where sr.form = p_form_id and sr."user" = p_user_id;

    if v_form.opens > now() then
        raise exception 'Form not yet open.' using hint = 'form-closed';
    end if;

    if v_form.closes < now() then
        raise exception 'Form closed.' using hint = 'form-closed';
    end if;

    if v_form.max_responses is not null and v_form.max_responses <= v_total_count then
        raise exception 'Form reached maximum responses.' using hint = 'form-closed';
    end if;

    if v_existing_count is not null and v_form.individual_limit <= v_existing_count then
        raise exception 'Maximum responses submitted.' using hint = 'form-closed';
    end if;

    insert into responses (form, respondent) values (p_form_id, p_user_id)
    returning * into v_response;

    insert into submission_records (form, "user", responses)
    values (p_form_id, p_user_id, 1) on conflict (form, "user") do update
    set responses = submission_records.responses + excluded.responses;

    return v_response;
end;
$$ language plpgsql;

create or replace function get_response_by_id(
    p_id text,
    p_form_id text,
    p_user_id text
) returns responses as $$
declare
    v_response responses;
begin
    select * into v_response from responses where id = p_id;
    if not found then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, p_form_id, 'analyze'::permission_role) then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    elsif v_response.respondent != p_user_id then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    return v_response;
end;
$$ language plpgsql;

create or replace function get_answers_for_response(
    p_id text,
    p_form_id text,
    p_user_id text
) returns setof answers as $$
declare
    v_respondent text;
begin
    select respondent into v_respondent from responses r where r.id = p_id;
    if not found then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, p_form_id, 'analyze'::permission_role) then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    elsif v_respondent != p_user_id then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    return query select * from answers a where a.response = p_id; 
end;
$$ language plpgsql;

create or replace function add_answer_to_response(
    p_id text,
    p_form_id text,
    p_user_id text,
    p_question text,
    p_value text
) returns answers as $$
declare
    v_answer answers;
    v_respondent text;
begin
    select respondent into v_respondent from responses r where r.id = p_id;
    if not found then
        raise exception 'Response not found or you do not have permission do this 1.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, p_form_id, 'respond'::permission_role) then
        raise exception 'Response not found or you do not have permission do this 2.' using hint = 'forbidden';
    end if;

    if v_respondent is not null and v_respondent != p_user_id then
        raise exception 'Response not found or you do not have permission do this 3.' using hint = 'forbidden';
    end if;

    insert into answers (response, question, value) values (
        p_id, p_question, p_value::jsonb
    ) on conflict (response, question) do update
    set value = excluded.value returning * into v_answer;

    return v_answer;
end;
$$ language plpgsql;

create or replace function submit_response_by_id(
    p_id text,
    p_form_id text,
    p_user_id text,
    p_save boolean
) returns responses as $$
declare
    v_response responses;
begin
    select * into v_response from responses r where r.id = p_id;
    if not found then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if not has_form_permission(p_user_id, p_form_id, 'respond'::permission_role) then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if v_response.respondent is not null and v_response.respondent != p_user_id then
        raise exception 'Response not found or you do not have permission do this.' using hint = 'forbidden';
    end if;

    if v_response.submitted is not null then
        update responses set status = 'edited', edited = now()
        where responses.id = p_id returning * into v_response;
    else
        update responses set status = 'completed', submitted = now()
        where responses.id = p_id returning * into v_response;
    end if;

    if p_save then
        insert into saved_responses ("user", form, response) values (
            p_user_id, p_form_id, p_id
        ) on conflict do nothing;
    end if;

    return v_response;
end;
$$ language plpgsql;

create or replace function list_responses_for_user(
    p_user_id text,
    p_form_title text,
    p_status response_status,
    p_sort_by text,
    p_order text,
    p_limit int,
    p_offset int
) returns setof responses as $$
begin
    return query select r.* from saved_responses sr
    join responses r on r.id = sr.response join forms f on f.id = sr.form
    where sr."user" = p_user_id and (p_status is null or r.status = p_status)
    and (p_form_title = '' or f.title %> p_form_title) order by
        similarity(f.title, p_form_title) desc,
        case when p_sort_by = 'submitted' and p_order = 'asc' then r.submitted end asc,
        case when p_sort_by = 'submitted' and p_order = 'desc' then r.submitted end desc,
        case when p_sort_by = 'started' and p_order = 'asc' then r.started end asc,
        case when p_sort_by = 'started' and p_order = 'desc' then r.started end desc,
        case when p_sort_by = 'edited' and p_order = 'asc' then r.edited end asc,
        case when p_sort_by = 'edited' and p_order = 'desc' then r.edited end desc
    limit p_limit offset p_offset;
end;
$$ language plpgsql;

create or replace function count_responses_for_user(
    p_user_id text,
    p_form_title text,
    p_status response_status
) returns bigint as $$
declare
    v_count bigint;
begin
    select count(distinct r.id) into v_count from saved_responses sr
    join responses r on r.id = sr.response join forms f on f.id = sr.form
    where sr."user" = p_user_id and (p_status is null or r.status = p_status)
    and (p_form_title = '' or f.title %> p_form_title);

    return v_count;
end;
$$ language plpgsql;
