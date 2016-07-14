package erroneous

import (
	"fmt"
	"sync"
	"time"
)

// MemoryStore implements a simple in-memory store for errors. The zero-value MemoryStore{} should be ready to use immediately.
type MemoryStore struct {
	errors []*Error
	sync.Mutex
	RollupDuration time.Duration
}

// ProtectError marks an error as protected
func (m *MemoryStore) ProtectError(guid Guid) error {
	m.Lock()
	defer m.Unlock()
	for _, e := range m.errors {
		if e.GUID.String() == guid.String() {
			e.Protected = true
			return nil
		}
	}
	return fmt.Errorf("Error %s not found", guid)
}

func (m *MemoryStore) DeleteError(guid Guid) error {
	m.Lock()
	defer m.Unlock()
	for i, e := range m.errors {
		if e.GUID.String() == guid.String() {
			m.removeAt(i)
			return nil
		}
	}
	return fmt.Errorf("Error %s not found", guid)
}

//remove an element.
func (m *MemoryStore) removeAt(i int) {
	if i >= len(m.errors) {
		return
	}
	a := m.errors
	m.errors = append(a[:i], a[i+1:]...)
}

func (m *MemoryStore) DeleteAllErrors(appName string) error {
	m.Lock()
	defer m.Unlock()
	//iterate backward so we don't mess up prior indexes
	for i := len(m.errors) - 1; i >= 0; i-- {
		err := m.errors[i]
		if err.Protected {
			continue
		}
		if appName == "" || err.ApplicationName == appName {
			m.removeAt(i)
		}
	}
	return nil
}

func (m *MemoryStore) LogError(e *Error) error {
	var dup *Error
	m.Lock()
	defer m.Unlock()
	if m.RollupDuration != 0 {
		minTime := e.CreationDate.Add(-1 * m.RollupDuration)
		for _, err := range m.errors {
			if e.ErrorKey == err.ErrorKey && err.CreationDate.After(minTime) {
				dup = err
				break
			}
		}
	}
	if dup != nil {
		e.GUID = dup.GUID
		dup.Count++
	} else {
		m.errors = append(m.errors, e)
	}
	return nil
}
func (m *MemoryStore) GetError(guid Guid) (*Error, error) {
	fmt.Printf("LOOKING FOR %s\n", guid)
	m.Lock()
	defer m.Unlock()
	for _, e := range m.errors {
		if e.GUID.String() == guid.String() {
			return e, nil
		}
	}
	return nil, fmt.Errorf("Error %s not found", guid)
}

func (m *MemoryStore) GetAllErrors(appName string) ([]*Error, error) {
	m.Lock()
	defer m.Unlock()
	matches := make([]*Error, 0, len(m.errors))
	for _, e := range m.errors {
		if appName == "" || e.ApplicationName == appName {
			matches = append(matches, e)
		}
	}
	return matches, nil
}
func (m *MemoryStore) GetErrorCount(appName string, since time.Time) (int, error) {
	m.Lock()
	defer m.Unlock()
	count := 0
	for _, e := range m.errors {
		if appName == "" || e.ApplicationName == appName {
			if e.CreationDate.After(since) {
				count++
			}
		}
	}
	return count, nil
}
