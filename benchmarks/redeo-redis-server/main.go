package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/bsm/redeo"
	"github.com/bsm/redeo/resp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var port int
	flag.IntVar(&port, "port", 6380, "server port")
	flag.Parse()

	go log.Printf("started server at :%d", port)
	srv := redeo.NewServer(nil)
	srv.HandleFunc("ping", func(w resp.ResponseWriter, _ *resp.Command) {
		w.AppendInlineString("PONG")
	})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer lis.Close()

	return srv.Serve(lis)
}
