#!/usr/bin/env fish

docker run --name fdb -dp 5432:5432 \
    -e POSTGRES_DB=forms -e POSTGRES_USER=super \
    -e POSTGRES_PASSWORD=super-duper-pizza-lining \
    postgres:alpine
