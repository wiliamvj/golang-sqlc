package main

import (
  "context"
  "database/sql"
  "fmt"
  "log"
  "os"

  "github.com/joho/godotenv"
  _ "github.com/lib/pq"
  "github.com/wiliamvj/golang-sqlc/internal/database/db"
  "github.com/wiliamvj/golang-sqlc/internal/handler"
)

func main() {
  // load .env file
  godotenv.Load()

  postgresURI := os.Getenv("DATABASE_URL")
  dbConnection, err := sql.Open("postgres", postgresURI)
  if err != nil {
    log.Panic(err)
  }
  err = dbConnection.Ping()
  if err != nil {
    dbConnection.Close()
    log.Panic(err)
  }

  fmt.Println("Connected to database")

  // start slqc queries
  slqcQueries := db.New(dbConnection)
  q := handler.UserHandler{
    Queries: slqcQueries,
    Db:      dbConnection,
  }
  handler.CreateUser(context.Background(), &q)

  // keep the program running
  select {}
}
