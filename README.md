# pr-compliance-app

This GitHub App create a status check on all PRs that require a form to be submitted before it can be merged.

# Requirements

* The app uses [smee.io](https://smee.io/) and docker-compose to run the app  and database (containers) locally. 
* You must create a github app and point the `Webhook URL` to your smee.io url.
* You must then install that github app on your account and select the repo(s) you would like to use it.
* Note that the status checks created on your PR will point to `http://localhost:8080`
	* This is hard-coded for simplicity, but would need to be changed in a real-world deployment.
	* Make sure you don't have anything else listening on that port before running the app locally.

## Environment variables

You must create a `.env` file that docker-compose will use to run your app. This will be populated with the private key and client secret that was used to create your github app.

### Example

Note that the `DATABASE_URL` should not be changed when testing the app locally with docker-compose.

```bash
DATABASE_URL=postgres://postgres:postgres@postgres:5432/pr_compliance
GITHUB_APP_IDENTIFIER=123456
GITHUB_WEBHOOK_SECRET=some secret data
GITHUB_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----
private key here ...
-----END RSA PRIVATE KEY-----"
```


# Running the app (locally)

* You will need to build the binary and create a container for the app. We use [paketo](https://paketo.io/) [cloud native buildpacks](https://buildpacks.io/) to do that with one command.
* You then use docker-compose to bring the app up or down.

## Commands


```bash
# Setup (these commands only need to be run once)
pack config default-builder paketobuildpacks/builder:base
git clone https://github.com/jericop/pr-compliance-app
cd pr-compliance-app
echo "create your .env file in this directory"

# Build the container and run the db and app with docker-compose
# This can be run repeatedly, such as whenever you are ready to test a code change.
pack build pr-compliance-app && docker-compose down -v --remove-orphans && docker-compose up -d

# Set the smee.io channel to point to your app listening on localhost:8080
# This is a blocking command that you may want to run in another window/tab.
# If the app stops receiving webhook events restarting this should fix it.
smee --url https://smee.io/some-unique-id --path /webhook_events --port 8080
```

# App Details

* A postgres database is as the backend.
* The `db/init` and `db/migrations` folders contain the DDL for app (users, database, tables).
* Changes are applied using the golang tool [migrate](https://github.com/golang-migrate/migrate)
* The `queries` folder contains queries used by the application.
	* They are annotated with comments for sqlc.
* [sqlc](https://docs.sqlc.dev/en/latest/) is used to generate a `Querier` interface and struct that implements all the annotated queries found in the `queries` folder.
* The `Querier` interface is used to interact with postgres.
* [faux](https://pkg.go.dev/github.com/ryanmoran/faux) is used to generate a mock of the `Querier` interface which is used in unit tests.
* [test2doc](https://pkg.go.dev/github.com/s-mang/test2doc) is used to generate an [API Blueprint](https://github.com/apiaryio/api-blueprint/blob/master/API%20Blueprint%20Specification.md) document from unit tests.
	* I learned about this project on the gotime fm podcast [here](https://changelog.com/gotime/5).
	* This dependency is only used in test code.

# Tests

This project strives for high test coverage, especially where it matters most. Since `test2doc` is  used to generate the API Blueprint document, all of the api endpoints should be tested. This also assumes that any api changes will come with full tests, so the API Blueprint document stays up to date.

# TODO

* Use database transactions for `/webhook_events` endpoint because multiple tables are modified before the status check is created.