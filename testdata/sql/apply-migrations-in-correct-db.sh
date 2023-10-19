#!/bin/bash -x

APP_DB="pr_compliance"
APP_MIG_USER_NAME="pr_compliance_migrations"
APP_MIG_USER_PASS="pr_compliance_migrations"

export PGPASSWORD="$APP_MIG_USER_PASS" # avoids password prompt

set -x

for f in $(ls /docker-entrypoint-initdb.d/*.sql | sort | grep -v 101.sql); do   
  psql -v ON_ERROR_STOP=1 --username $APP_MIG_USER_NAME --dbname $APP_DB < ${f}
done
