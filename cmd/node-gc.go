package main

import (
	"time"

	"github.com/qiujian16/node-gc/pkg/controller"
)

func main() {
	period := 5 * time.Second

	stopCh := make(chan struct{})

	gcController, err := controller.NewGCController(period)
	if err != nil {
		panic(err)
	}

	go gcController.Run()
	<-stopCh
}
