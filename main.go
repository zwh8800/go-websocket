package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"fmt"
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
					fmt.Printf("client [%s] closed\n", ws.LocalAddr())
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
				fmt.Printf("received: %s\n", receive)

				sendMsg := Payload{
					"message",
					"[" + receive.Msg + "]",
				}
				err := websocket.JSON.Send(ws, sendMsg)
				if err != nil {
					log.Fatal("Write: ", err)
				}
				fmt.Printf("sent: %s\n", sendMsg)
			case <- closeCh:
				return
			}
		}
	}()

	<- closeCh
	ws.Close()
}

func main() {
	ch1 = make(chan Payload)

	http.Handle("/", http.FileServer(http.Dir("public/dist/")))
	http.Handle("/ws", websocket.Handler(handler))
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
