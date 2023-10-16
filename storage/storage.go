package storage

//SKIP go:generate mockgen -package testmocks -destination mock/querier.go github.com/jericop/pr-compliance-app/storage/postgres Querier
//go:generate faux --interface Querier --output ../fakes/querier_mock.go --package github.com/jericop/pr-compliance-app/storage/postgres
