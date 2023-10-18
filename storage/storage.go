package storage

// This will generate the fakes.Querier struct that implements the postgres.Querier interface.
//go:generate faux --interface Querier --output ../fakes/querier.go --package github.com/jericop/pr-compliance-app/storage/postgres
