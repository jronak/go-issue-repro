package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func run() {
	go runServer()
	runClient()
}

func runServer() {
	h2server := &http2.Server{}
	gServer := grpc.NewServer(grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
		msg := []byte{}
		if err := stream.RecvMsg(&msg); err != nil {
			return err
		}
		return stream.SendMsg(msg)
	}))
	server := &http.Server{
		Handler: h2c.NewHandler(gServer, h2server),
		Addr:    ":8081",
	}

	fmt.Printf("H2c Server starting on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func runClient() {
	payload := generatePayload(1 << 16) // 16KB payload.
	conn, err := grpc.Dial("localhost:8081", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "content-length", "2")
	var res []byte
	for i := 0; ; i++ {
		err = conn.Invoke(ctx, "/golang.org.Repro/Issue", payload, &res, grpc.ForceCodec(codec{}))
		if err != nil {
			log.Printf("Failed stream num:%d with err:%v\n", i, err)
		} else {
			log.Printf("Successful stream num: %d with resp:%v\n", i, res)
		}
	}
}

func generatePayload(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = 'A'
	}

	return b
}
