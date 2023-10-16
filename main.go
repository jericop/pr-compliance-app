package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jericop/pr-compliance-app/api"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

func main() {
	fmt.Printf("DATABASE_URL=%s\n", os.Getenv("DATABASE_URL"))

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	connPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer connPool.Close()

	db := postgres.New(connPool)
	apiServer, err := api.NewServer(connPool, db)
	if err != nil {
		log.Fatalf("Error setting up server: %v\n", err)
	}

	http.ListenAndServe(":8080", apiServer.Router)
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
