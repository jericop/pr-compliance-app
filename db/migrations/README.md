# database migrations

The golang tool [migrate](https://github.com/golang-migrate/migrate) will be used to apply database migrations.

The examples given below were tested using an instance of postgres running locally with the username and password as `postgres`.

It is important to test that migrations work for both up and down. It is also recommended that up/down migrations be applied immediately to verify there are no issues.

Applying all migrations 
```
migrate -database postgres://postgres:postgres@localhost:5432/pr_compliance\?sslmode=disable -path db/migrations up
```

Applying only the latest migrations
```
migrate -database postgres://postgres:postgres@localhost:5432/pr_compliance\?sslmode=disable -path db/migrations up 1
```

Removing all migrations 
```
migrate -database postgres://postgres:postgres@localhost:5432/pr_compliance\?sslmode=disable -path db/migrations down --all
```