# DB migration (MySQL)

## how to

```bash
./tools/bin/migrate create -ext sql -dir db/migrations -seq <do_something>
```

### how to locally practice

using `mycli`

```mycli
source db/migrations/000001_create_users.up.sql;
```

```
./tools/bin/migrate -database "postgres://see:seedbpass@see-db:54321/see?sslmode=disable" -path db/migrations up
```