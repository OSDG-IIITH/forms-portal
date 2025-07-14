create table if not exists users (
    id text primary key default generate_ulid(),
    handle text not null unique, -- cas user id
    email text not null unique,
    name text not null
);

create table if not exists forms (
    id text primary key default generate_ulid(),
    owner text not null references users(id) on delete cascade,
    slug text not null,
    title text not null,
    description text,
    structure text not null,
    modified timestamptz not null default now(),
    live boolean not null default false,
    opens timestamptz,
    closes timestamptz,
    max_responses int,
    individual_limit int not null default 1,
    editable_responses boolean not null default false,

    unique (owner, slug),
    constraint response_limits_check check (
        individual_limit >= 1 and max_responses > individual_limit
    )
);

create table if not exists groups (
    id text primary key default generate_ulid(),
    owner text not null references users(id),
    name text not null,
    description text,
    type group_type not null,

    unique (owner, name)
);

-- note: this table is empty, only exists for sqlc to understand the type
create table if not exists group_with_details (
    id text, owner text, name text, description text,
    type group_type not null, domain text, members text[]
);

create table if not exists group_domain_rules (
    "group" text primary key references groups(id) on delete cascade,
    domain text not null unique
);

create table if not exists group_list_members (
    "group" text not null references groups(id) on delete cascade,
    "user" text not null references users(id) on delete cascade,

    primary key ("group", "user")
);

create table if not exists form_permissions (
    id text primary key default generate_ulid(),
    form text not null references forms(id) on delete cascade,
    role permission_role not null,
    "user" text references users(id) on delete cascade,
    "group" text references groups(id) on delete cascade,

    constraint permit_user_or_group check (
        ("user" is not null and "group" is null) or
        ("user" is null and "group" is not null)
    )
);

create unique index form_permissions_user_unique
on form_permissions (form, "user", role) where "group" is null;

create unique index form_permissions_group_unique
on form_permissions (form, "group", role) where "user" is null;

create table if not exists submission_records (
    form text not null references forms(id) on delete cascade,
    "user" text not null references users(id) on delete cascade,
    responses int not null default 1,

    primary key ("form", "user")
);

create table if not exists responses (
    id text primary key,
    form text not null references forms(id) on delete cascade,
    respondent text references users(id),
    status response_status not null default 'draft',
    started timestamptz not null default now(),
    submitted timestamptz,
    edited timestamptz
);

create table if not exists answers (
    id text primary key default generate_ulid(),
    response text not null references responses(id) on delete cascade,
    question text not null, -- id of question in forms.spec
    value jsonb not null,
    submitted timestamptz not null default now(),
    modified timestamptz not null default now(),

    unique (response, question)
);

create table if not exists comments (
    id text primary key default generate_ulid(),
    form text not null references forms(id) on delete cascade,
    commenter text not null references users(id) on delete cascade,
    body text not null,
    state comment_state not null default 'visible',
    element text not null, -- id of question/section in forms.spec
    parent text references comments(id) on delete cascade,
    modified timestamptz not null default now()
);
