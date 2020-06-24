package monitor

import (
	"log"
	"net"
	"sync"
	"time"
)

type Check struct {
	host string
	port string
	open bool
}

type Storage interface {
	SaveCheck(Check) error
	GetChecks() ([]Check, error)
}

type Monitor struct {
	storage Storage
	logger  *log.Logger
	wg      *sync.WaitGroup
	done    chan struct{}
}

func NewMonitor(storage Storage, logger *log.Logger) *Monitor {
	return &Monitor{
		storage: storage,
		logger:  logger,
		wg:      &sync.WaitGroup{},
		done:    make(chan struct{}),
	}
}

func (m *Monitor) Start() {
	m.wg.Add(1)
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-m.done:
			m.logger.Println("Monitor done")
			m.wg.Done()
			return
		case <-ticker.C:
			err := m.runChecks()
			if err != nil {
				m.logger.Println(err)
			}
		}

	}

}

func (m *Monitor) Stop() {
	close(m.done)
	m.wg.Wait()
}

func (m *Monitor) AddCheck(host string, port string, open bool) error {
	check := Check{
		host,
		port,
		open,
	}

	err := m.storage.SaveCheck(check)

	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) runChecks() error {
	checks, err := m.storage.GetChecks()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(checks))

	for _, check := range checks {
		go m.runCheck(&check, &wg)
	}

	wg.Wait()

	return nil

}

func (m *Monitor) runCheck(check *Check, wg *sync.WaitGroup) {
	defer wg.Done()

	open, err := portOpen(check)

	if err != nil {
		if _, ok := err.(net.Error); !ok {
			m.logger.Println(err)
			return
		}
	}

	if open != check.open {
		var expected string
		if check.open {
			expected = "open"
		} else {
			expected = "closed"
		}
		m.logger.Printf("%s:%s should be %s", check.host, check.port, expected)
	}
}

func portOpen(check *Check) (bool, error) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(check.host, check.port), timeout)

	if err != nil {
		return false, err
	}

	if conn != nil {
		defer conn.Close()
		return true, nil
	}

	return false, nil
}
