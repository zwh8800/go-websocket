package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"fmt"
	"time"
	"math/rand"
	"strconv"
	"io"
)

type Payload struct {
	Type string `json:"type"`
	Msg string	`json:"msg"`
}

var ch1 chan Payload

func handler(ws *websocket.Conn) {
	closeCh := make(chan int)

	go func() {
		var receive Payload
		for {
			err := websocket.JSON.Receive(ws, &receive)
			if err != nil {
				if (err.Error() == "EOF") {
					fmt.Printf("client [%#v] closed\n", ws.LocalAddr())
					close(closeCh)	//close connection
					return
				} else {
					log.Fatal("Read: ", err)
				}
			}
			ch1 <- receive
		}
	}()

	go func() {
		for {
			select {
			case receive := <- ch1:
				fmt.Printf("received: %#v\n", receive)

				sendMsg := Payload{
					"message",
					"[" + receive.Msg + "]",
				}
				err := websocket.JSON.Send(ws, sendMsg)
				if err != nil {
					log.Fatal("Write: ", err)
				}
				fmt.Printf("sent: %#v\n", sendMsg)
			case <- closeCh:
				return
			}
		}
	}()

	<- closeCh
	ws.Close()
}



type Message struct {
	Type string `json:"type"`
	Token int `json:"token"`
	Data interface{} `json:"data"`
}
func handler2(ws *websocket.Conn) {
	for {
		var receive Message
		err := websocket.JSON.Receive(ws, &receive)
		if err != nil {
			if (err == io.EOF) {
				fmt.Printf("client [%#v] closed\n", ws.LocalAddr())
				return
			} else {
				log.Fatal("Read: ", err)
			}
		}
		fmt.Printf("received: %#v\n", receive)

		go func() {
			switch receive.Type {
			case "GetChannels":
				timeout := 100 + rand.Intn(200);
				time.Sleep(time.Duration(timeout) * time.Millisecond)
				sendMsg := Message{
					"GetChannelsResponse",
					receive.Token,
					timeout,
				}
				err := websocket.JSON.Send(ws, sendMsg)
				if err != nil {
					log.Fatal("Write: ", err)
				}

			case "ListChannels":
				timeout := 300 + rand.Intn(200);
				time.Sleep(time.Duration(timeout) * time.Millisecond)
				sendMsg := Message{
					"ListChannelsResponse",
					receive.Token,
					strconv.Itoa(timeout),
				}
				err := websocket.JSON.Send(ws, sendMsg)
				if err != nil {
					log.Fatal("Write: ", err)
				}

			case "Join":
				timeout := 300 + rand.Intn(200);
				time.Sleep(time.Duration(timeout) * time.Millisecond)
				sendMsg := Message{
					"JoinResponse",
					receive.Token,
					map[string] interface{}{
						"time": timeout,
						"channel": receive.Data,
					},
				}
				err := websocket.JSON.Send(ws, sendMsg)
				if err != nil {
					log.Fatal("Write: ", err)
				}
			}
		}()
	}

}

func main() {
	ch1 = make(chan Payload)

	http.Handle("/", http.FileServer(http.Dir("public/dist/")))
	http.Handle("/ws", websocket.Handler(handler2))
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
