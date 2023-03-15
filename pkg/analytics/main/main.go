package main

import (
	"fmt"
	"github.com/rudderlabs/analytics-go/v4"
)

func main() {
	// Instantiates a client to use send messages to the Rudder API.
	client := analytics.New("2N28inls9mXwq12IUo0yB34rFsm", "https://selefralefsm.dataplane.rudderstack.com")

	// Enqueues a track event that will be sent asynchronously.
	err := client.Enqueue(analytics.Track{
		UserId: "test-user",
		Event:  "test-snippet",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Flushes any queued messages and closes the client.
	err = client.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("send success")
}
