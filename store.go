package erroneous

import (
	"time"

	"github.com/twinj/uuid"
)

// Guid is an alias for uuid.UUID. Indluded to reduce repetition of external package
type Guid uuid.UUID

// Store is an interface to abstract storage of error data.
type Store interface {
	ProtectError(Guid) error
	DeleteError(Guid) error
	DeleteAllErrors(appName string) error
	LogError(*Error) error
	GetError(Guid) (*Error, error)
	GetAllErrors(appName string) ([]*Error, error)
	GetErrorCount(appName string, since time.Time) (int, error)
}
