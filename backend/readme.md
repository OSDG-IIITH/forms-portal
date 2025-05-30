## Forms Portal Backend

> This directory contains the source code for the backend of the Forms Portal.

The backend is written in `go`, using the `echo` server framework and `postgres`
as the database service.

This project uses [`mise`](https://github.com/jdx/mise) for managing tools.
For running the database locally, `docker` is used.

Before starting the server, copy the `.env.example` file to `.env`:

```
cp .env.example .env
```

The server can be started by running:

```fish
go run .
```

The database can be setup by running:

```fish
./database/scripts/create.fish
./database/scripts/setup.fish
```

If you are not using the `fish` shell, view the scripts and run the commands
yourself using your shell's syntax.

To open up a `psql` shell, run:

```fish
docker exec -it fdb -U field --dbname forms
```
