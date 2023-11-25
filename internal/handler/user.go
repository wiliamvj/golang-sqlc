package handler

import (
  "context"
  "database/sql"
  "fmt"
  "time"

  "github.com/google/uuid"
  "github.com/wiliamvj/golang-sqlc/internal/database/db"
)

type UserHandler struct {
  Queries *db.Queries
  Db      *sql.DB
}

func runWithTX(ctx context.Context, c *sql.DB, fn func(*db.Queries) error) error {
  tx, err := c.BeginTx(ctx, nil)
  if err != nil {
    return err
  }
  q := db.New(tx)
  err = fn(q)
  if err != nil {
    if errRb := tx.Rollback(); errRb != nil {
      return fmt.Errorf("error on rollback: %v, original error: %w", errRb, err)
    }
    return err
  }
  return tx.Commit()
}

func CreateUser(ctx context.Context, h *UserHandler) error {
  userID := uuid.New().String()
  userEmail := fmt.Sprintf("john.doe-%v@email.com", time.Now().Unix())

  err := runWithTX(ctx, h.Db, func(q *db.Queries) error {
    var err error
    err = q.CreateOneUser(ctx, db.CreateOneUserParams{
      ID:        userID,
      Name:      "John Doe",
      Email:     userEmail,
      Password:  "123456",
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
    })
    if err != nil {
      fmt.Println("Error creating user", err)
      return err
    }
    fmt.Println("User created")

    // create post
    err = q.CreateOnePost(ctx, db.CreateOnePostParams{
      ID:        uuid.New().String(),
      Title:     "SLQC with Golang",
      Body:      "This is a post about SLQC with Golang",
      AuthorID:  userID,
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
    })
    if err != nil {
      fmt.Println("Error creating post", err)
      return err
    }
    fmt.Println("Post created")
    return nil
  })
  if err != nil {
    fmt.Println("Error creating user and post, roll back applied", err)
    return err
  }
  return nil
}
