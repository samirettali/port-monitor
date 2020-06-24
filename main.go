package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/samirettali/port-monitor/monitor"
	"github.com/samirettali/port-monitor/storage"
)

func main() {
	logger := log.New(os.Stdout, "", 0)
	storage := storage.NewInMemoryStorage()

	logger.Println("Starting monitor")
	monitor := monitor.NewMonitor(storage, logger)

	// Just for testing
	for i := 1; i <= 65535; i++ {
		n := strconv.Itoa(i)
		monitor.AddCheck("localhost", n, false)
	}

	go monitor.Start()

	ticker := time.NewTicker(time.Second * 20)
	<-ticker.C
	monitor.Stop()
	logger.Println("Exiting main")
}
