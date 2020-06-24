package storage

import (
	"github.com/samirettali/port-monitor/monitor"
)

type InMemoryStorage struct {
	checks []monitor.Check
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		checks: make([]monitor.Check, 0),
	}

}

func (s *InMemoryStorage) SaveCheck(check monitor.Check) error {
	s.checks = append(s.checks, check)
	return nil
}

func (s *InMemoryStorage) GetChecks() ([]monitor.Check, error) {
	return s.checks, nil
}
