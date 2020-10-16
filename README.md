# Godbmigrate

GoLang database migrations made simple.

## Installation

The classic way of installing it using go get command

    go get -u github.com/josemiguelmelo/godbmigrate

## Usage

### Configuration file

First of all, you need to add a configuration file to your project with the database connection information. By default, this file must be located on the root of your project and should be named `godbmigrate.yml`.

You can specify another configuration file location, by passing it as a command parameter:

    godbmigrate -migrate -config-location=./configs/database

[Example of a configuration file:](./example/godbmigrate.yml)

```yaml
database:
    provider: postgres
    connectionUri: "user=example password=password DB.name=exampledb port=5432 host=localhost sslmode=disable"
    # if connectionUri is set, don't need to add the following:
    username: example
    password: password
    database: exampledb
    host: localhost
    port: 5432
```

### Migration file

The migration file should be an SQL file and should be placed inside `db/migrations` folder.

A migration is composed by two sections: 

-   `--- up! ---` - used to apply new migrations
-   `--- down! ---` - used for rolling back migrations

[Example:](./example/db/migrations/1.sql)

```sql
--- up! ---

CREATE TABLE example(
    example_col VARCHAR(240)
)

--- down! ---

DROP TABLE IF EXISTS example;
```

### Run migration

Run the following command to run migrations on the root of your project:

    godbmigrate -migrate

### Rollback migration

Run the following command to rollback previously runned migrations:

    godbmigrate -migrate-rollback

## License

[MIT License](LICENSE)
