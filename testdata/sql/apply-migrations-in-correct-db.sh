#!/bin/bash -x

for f in $(ls /docker-entrypoint-initdb.d/*.sql | sort | grep -v 101.sql); do 
  echo "\$f=$f"
  # psql -U postgres -h localhost -d pr_compliance < ${f}
  psql -U postgres -d pr_compliance < ${f}
done
