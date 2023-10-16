-- This only needs to be run once when the database is created.

-- Workaround for `CREATE DATABASE IF NOT EXISTS`
SELECT 'CREATE DATABASE pr_compliance'
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'pr_compliance'
) \gexec


SELECT 'CREATE USER pr_compliance_api'
WHERE NOT EXISTS (
    SELECT FROM pg_catalog.pg_roles WHERE rolname = 'pr_compliance_api'
) \gexec
GRANT CONNECT ON DATABASE pr_compliance TO pr_compliance_api;


SELECT 'CREATE USER pr_compliance_migrations'
WHERE NOT EXISTS (
    SELECT FROM pg_catalog.pg_roles WHERE rolname = 'pr_compliance_migrations'
) \gexec
GRANT CONNECT ON DATABASE pr_compliance TO pr_compliance_migrations;
