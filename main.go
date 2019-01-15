package main

import (
	"github.com/IvoryRaptor/iotbox/kernel"
)

func main() {
	if box, err := kernel.Boot(); err == nil {
		box.Start()
	} else {
		println(err.Error())
	}
	select {}
}
