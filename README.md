# Add/Drop Classes
Basic implementation for an Add/Drop Classes web application

```
Usage: add-drop-classes-api <command>

A WebApp for Add/Drop Classes in a college

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  server
    Run the server

Run "add-drop-classes-api <command> --help" for more information on a command.
```

## Development

Since this application uses `postgreSQL`, there is a `docker-compose.yaml` file available in the project root.

Running `docker-compose up -d` will run a `postgreSQL` and `pgadmin4` instance.

Going to http://127.0.0.1:9000 on the browser will open the `pgadmin4` where logging in with the following credentials
are possible;
```
email: admin@admin.com 
password: admin
```

After logging in, it is possible to add `postgreSQL` database as a server to the `pgadmin4`. Since the `postgreSQL`
runs in the same docker network with the `pgadmin4`, the host name for the `postgreSQL` instance is `adc_postgres`.
Username and password are both `postgres`.

Since `schema.sql` file is added as a volume to the `postgreSQL` container, it is automatically executed when the container
is being run.

It's possible to follow the logs outputted by the containers by `docker-compose logs -f`.

Running `docker-compose down` will stop all the containers. If there is a need to reset the data in the database deleting
the directory `./docker/data/postgres` would suffice.
