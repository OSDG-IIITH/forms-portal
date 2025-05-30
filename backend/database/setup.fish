#!/bin/fish

docker cp database/ fdb:/database/
for f in database/schemas/*.sql
    docker exec -it fdb psql -U super -d forms -f /$f
end
