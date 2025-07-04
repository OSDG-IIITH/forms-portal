create table if not exists users (
    id text primary key default generate_ulid(),
    handle text not null, -- cas user id
    name text not null,
    email text not null unique
);

create table if not exists forms (
    id text primary key default generate_ulid(),
    owner text not null references users(id) on delete cascade,
    title text not null,
    slug text not null,
    description text,
    specification jsonb not null,
    salt uuid not null default gen_random_uuid(),
    modified timestamptz not null default now(),
    live boolean not null default false,
    opens timestamptz,
    closes timestamptz,
    unique (owner, slug)
);

create table if not exists groups (
    id text primary key default generate_ulid(),
    owner text not null references users(id),
    name text not null,
    description text,
    type group_type not null,
    unique (owner, name)
);

create table if not exists group_domain_rules (
    group_id text primary key references groups(id) on delete cascade,
    domain text not null unique
);

create table if not exists group_list_members (
    group_id text not null references groups(id) on delete cascade,
    user_id text not null references users(id) on delete cascade,
    primary key (group_id, user_id)
);

create table if not exists form_permissions (
    id text primary key default generate_ulid(),
    form text not null references forms(id) on delete cascade,
    role permission_role not null,
    user int references users(id) on delete cascade,
    group int references groups(id) on delete cascade,

    constraint permit_user_or_group check (
        (user is not null and group is null) or
        (user is null and group is not null)
    )
);

create table if not exists responses (
    id text primary key, -- hashed value of (user.id + form.salt)
    form text not null references forms(id) on delete cascade,
    respondent text references users(id),
    status response_status not null default 'in_progress',
    started timestamptz not null default now(),
    submitted timestamptz
);

create table if not exists answers (
    id text primary key default generate_ulid(),
    response not null references responses(id) on delete cascade,
    question text not null, -- id of question in forms.spec
    value jsonb not null,
    submitted timestamptz not null default now(),
    modified timestamptz not null default now(),

    unique (response, question)
);

create table if not exists comments (
    id text primary key default generate_ulid(),
    form int not null references forms(id) on delete cascade,
    user int not null references users(id) on delete cascade,
    element text not null, -- id of question/section in forms.spec
    body text not null,
    parent int references comments(id) on delete cascade,
    modified timestamptz not null default now(),
);
