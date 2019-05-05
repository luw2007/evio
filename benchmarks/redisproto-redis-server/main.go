package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/secmask/go-redisproto"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {
	go func() {
		var port int
		flag.IntVar(&port, "port", 6380, "server port")
		flag.Parse()
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Panic("tcp start error")
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error on accept: ", err)
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				parser := redisproto.NewParser(conn)
				writer := redisproto.NewWriter(bufio.NewWriter(conn))
				var ew error
				for {
					cmd, err := parser.ReadCommand()
					if err != nil {
						_, ok := err.(*redisproto.ProtocolError)
						if ok {
							ew = writer.WriteError(err.Error())
						} else {
							break
						}
					} else {
						switch strings.ToLower(string(cmd.Get(0))) {
						case  "ping":
							conn.Write([]byte("+PONG\r\n"))
						}
					}
					if cmd.IsLast() {
						writer.Flush()
					}
					if ew != nil {
						log.Fatal(err)
						break
					}
				}
			}(conn)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
