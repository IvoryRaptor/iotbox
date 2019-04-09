package main

import (
	"github.com/IvoryRaptor/iotbox/akka"
	"github.com/IvoryRaptor/iotbox/store/sqlite"
	"time"
)

func main() {
	system := &akka.System{}
	system.Start()
	system.ActorOf(&sqlite.SqliteActor{}, "sqlite")
	//system.Ask(sqlite,akka.Message{"":""},)
	time.Sleep(10 * time.Second)
}
