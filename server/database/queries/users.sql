-- name: EnsureUser :one
with x as (
    insert into users (handle, email, name) values ($1, $2, $3)
    on conflict do nothing returning *
) select * from x union select * from users where email = $2;

-- name: GetUserById :one
select * from users where id = $1;