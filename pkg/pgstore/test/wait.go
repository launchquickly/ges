package test

import (
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
)

// NewPostgresWaitStrategy wait strategy for use when testing against Postgres container as it restarts on startup.
//
// Waits until the log message "database system is ready to accept connections" is logged twice
// or timeout threshold of 60 seconds is exceeded.
func NewPostgresWaitStrategy() *wait.LogStrategy {
	return wait.NewLogStrategy("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(60 * time.Second)
}
