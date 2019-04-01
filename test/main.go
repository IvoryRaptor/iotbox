package main

import (
	"github.com/IvoryRaptor/iotbox/test/akka"
	"time"
)

func main() {
	system := &akka.System{}
	system.Start()
	test2(system)
	time.Sleep(10 * time.Second)
}
