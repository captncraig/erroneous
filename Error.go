package erroneous

import (
	"time"
)

//ErrorKey is the comparable part of the error. The only portion relevant to deduplicating.
type ErrorKey struct {
	MachineName string
	Detail      string
}

// Error is our representaion of an error. It is wraps the originating error with additional metadata, deduplication, and so forth.
type Error struct {
	ErrorKey
	Id              int64
	GUID            Guid
	Protected       bool
	Message         string
	ApplicationName string
	Count           int
	CreationDate    time.Time
}
