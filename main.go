package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"fmt"
)

var ch1 chan[]byte

func handler(ws *websocket.Conn) {
	closeCh := make(chan int)

	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := ws.Read(buffer)
			if err != nil {
				if (err.Error() == "EOF") {
					fmt.Printf("client [%s] closed\n", ws.LocalAddr())
					close(closeCh)	//close connection
					return
				} else {
					log.Fatal("Read: ", err)
				}
			}
			ch1 <- buffer[:n]
		}
	}()

	go func() {
		for {
			select {
			case buffer := <- ch1:
				n := len(buffer)
				msg := string(buffer)
				fmt.Printf("%d byte received: %s\n", n, msg)

				sendMsg := "[" + msg + "]"
				m, err := ws.Write([]byte(sendMsg))
				if err != nil {
					// TODO 判断连接是否关闭
					log.Fatal("Write: ", err)
				}
				fmt.Printf("%d byte sent: %s\n", m, string(([]byte(sendMsg))[:m]))
			case <- closeCh:
				return
			}
		}
	}()

	<- closeCh
	ws.Close()
}

func main() {
	ch1 = make(chan[]byte)

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/ws", websocket.Handler(handler))
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
