begin; create type permission_role as enum ('respond', 'view', 'comment', 'analyze', 'edit', 'manage'); commit;
begin; create type group_type as enum ('list', 'domain'); commit;
begin; create type response_status as enum ('draft', 'completed', 'edited'); commit;
