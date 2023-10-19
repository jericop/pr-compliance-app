# database init

The following sql statements will set up the database for the app. The `init-db.sh` script is used for local testing and implements the statements below.

# Create the database and users

* Run as user `postgres`
* Run in database `postgres`

```sql
CREATE DATABASE pr_compliance;

CREATE USER pr_compliance_api WITH ENCRYPTED PASSWORD '<pr_compliance_api password>';
GRANT CONNECT ON DATABASE pr_compliance TO pr_compliance_api;

CREATE USER pr_compliance_migrations WITH ENCRYPTED PASSWORD '<pr_compliance_migrations password>';
GRANT CONNECT, CREATE, TEMPORARY ON DATABASE pr_compliance TO pr_compliance_migrations;
```

# Create the schema

* Run as user `postgres`
* Run in database `pr_compliance`

```sql
CREATE SCHEMA production;

GRANT USAGE ON SCHEMA production TO pr_compliance_api;
GRANT USAGE, CREATE ON SCHEMA production TO pr_compliance_migrations;
```

# Grant user privileges

* Run as user `pr_compliance_migrations`
* Run in database `pr_compliance`

```sql
ALTER DEFAULT PRIVILEGES IN SCHEMA production GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO pr_compliance_api;
ALTER DEFAULT PRIVILEGES IN SCHEMA production GRANT ALL PRIVILEGES ON TABLES TO pr_compliance_migrations; 

ALTER DEFAULT PRIVILEGES IN SCHEMA production GRANT USAGE ON SEQUENCES TO pr_compliance_api;
ALTER DEFAULT PRIVILEGES IN SCHEMA production GRANT ALL PRIVILEGES ON SEQUENCES TO pr_compliance_migrations;
```