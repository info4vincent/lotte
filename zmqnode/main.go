package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/info4vincent/lotte"
	zmq "github.com/pebbe/zmq4"
)

func playSpeech(rubensTextToSay string) {
	fullFileName := lotte.GetSpeech(rubensTextToSay)
	lotte.PlayOgg(filepath.Base(fullFileName))
}

func main() {
	const PtrSize = 32 << uintptr(^uintptr(0)>>63)
	fmt.Println(runtime.GOOS, runtime.GOARCH)
	fmt.Println(strconv.IntSize, PtrSize)

	rubensTextToSay := "Hallo. Ik ben Ruben. Ik kom nu online om te luisteren..."
	playSpeech(rubensTextToSay)

	fmt.Println("starting lotte speech as ruben....")
	fmt.Println("Connecting to eventsource server...")
	requester, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatal(err)
	}

	defer requester.Close()
	err = requester.Connect("tcp://localhost:5555")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Client bind *:5555 succesful")

	requester.Send("zmqspeech available..", 0)

	event, error := requester.Recv(zmq.DONTWAIT)

	if error != nil {
		log.Println("Not received any..")
	} else {
		log.Println(event)
	}

	//  Socket to talk to server
	fmt.Println("Collecting updates from weather server...")
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5556")

	//  Subscribe to Everything which a command to say things..
	filter := "Say:"
	subscriber.SetSubscribe(filter)

	// Wait for messages
	for {
		sayCommand, error := subscriber.Recv(0)

		if error != nil {
			continue
		}

		// if !strings.HasPrefix(sayCommand, "Say:") {
		// 	log.Println("Event does not containing 'Say:'")
		// 	continue
		// }

		// send reply back to client
		reply := "received say command"
		requester.Send(reply, 0)

		rubensTextToSay := strings.TrimPrefix(sayCommand, "Say:")
		playSpeech(rubensTextToSay)
	}
}
