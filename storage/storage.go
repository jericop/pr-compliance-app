package storage

// This will generate the fake.Querier struct that implements Querier interface.
//go:generate faux --interface Querier --output ../fakes/querier.go --package github.com/jericop/pr-compliance-app/storage/postgres
