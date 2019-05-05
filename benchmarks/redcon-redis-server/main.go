package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/redcon"
	"log"
	"os"
	"os/signal"
	"strings"
)

func main() {
	go func() {
		var port int
		flag.IntVar(&port, "port", 6380, "server port")
		flag.Parse()
		err := redcon.ListenAndServe(fmt.Sprintf(":%d", port),
			func(conn redcon.Conn, cmd redcon.Command) {
				switch strings.ToLower(string(cmd.Args[0])) {
				case "ping":
					conn.WriteString("PONG")
				}
			},
			func(conn redcon.Conn) bool {
				// use this function to accept or deny the connection.
				// log.Printf("accept: %s", conn.RemoteAddr())
				return true
			},
			func(conn redcon.Conn, err error) {
				// this is called when the connection has been closed
				// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
