package db

//SKIP go:generate faux --interface Querier --output ../../fakes/querier_mock.go
//SKIP go:generate mockgen -package testmocks -destination ../testmocks/mock_docker_client.go github.com/jericop/pr-compliance-app/db Querier
