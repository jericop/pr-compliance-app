package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jericop/pr-compliance-app/api"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

func main() {
	fmt.Printf("DATABASE_URL=%s\n", os.Getenv("DATABASE_URL"))

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	connPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	fmt.Printf("connPool: %v\n", connPool)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer connPool.Close()

	db := postgres.New(connPool)

	// Create a new server with routes configured
	server := api.NewServer(db)

	// This is a wrapper that calls http.ListenAndServe, which is a blocking call.
	log.Fatal(server.Start(":8080"))
}

/*
package main

import (
	"fmt"

	"github.com/google/uuid"
	// "github.com/jericop/pr-compliance-app/api"
)

func main() {
	// api.HelloWorld()
	fmt.Println("uuid.New().String()")
	fmt.Println(uuid.New().String())
}
*/
