-- these indexes speed up the `has_form_permission` function
create index on form_permissions (form, role, "user");
create index on group_list_members ("user");
create index on group_domain_rules (domain);

-- these indexes help speed up full text search
create index on forms using gin (title gin_trgm_ops);
