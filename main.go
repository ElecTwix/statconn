package main

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"os"

	"github.com/ElecTwix/statconn/ping"
	"github.com/ElecTwix/statconn/ui"
)

func main() {
	f, err := os.OpenFile("network.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	localPingChn, RemotePingChn := make(chan error), make(chan error)
	localPingStop, LocalRunFunc := ping.PingAdress("192.168.1.1", localPingChn)
	defer localPingStop()
	remotePingStop, RemoteRunFunc := ping.PingAdress("1.1.1.1", RemotePingChn)
	defer remotePingStop()

	go LocalRunFunc()
	go RemoteRunFunc()

	err = ui.CreateUI(localPingChn, RemotePingChn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("bye")
}
