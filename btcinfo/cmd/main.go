package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	server "node-pinger/btcinfo/cmd/service"
	"node-pinger/btcinfo/pkg/endpoint"
	"node-pinger/btcinfo/pkg/service"
)

func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")

	flag.Parse()
	ctx := context.Background()
	srv := service.NewService()
	errChan := make(chan error)

	endpoints := endpoint.Endpoints{
		SyncStatusEndpoint: endpoint.MakeSyncStatusEndpoint(srv),
		BlockTxsEndpoint:   endpoint.MakeBlockTxsEndpoint(srv),
	}

	// Stop the server when ctrl+c is pressed
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// HTTP Transport (listen for incoming requests)
	go func() {
		log.Println("btcinfo is listening on port:", *httpAddr)
		handler := server.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
