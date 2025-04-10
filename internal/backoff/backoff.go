package backoff

import (
	"time"
)

type Config struct {
	// At is the time when the first failure happened.
	At time.Time
	// Retries is the number of attempts to connect.
	Retries int
}
