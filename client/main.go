package main

import (
	"connectrpc-stream-example/gen/proto"
	"connectrpc-stream-example/gen/proto/protoconnect"
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
)

func main() {
	fmt.Println("Hello")
	httpClient := http.DefaultClient
	url := "http://localhost:3000"
	client := protoconnect.NewMyProtoClient(httpClient, url)

	res, err := client.GiveInfo(context.TODO(), &connect.Request[proto.GiveInfoRequest]{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Msg.Data)

	stream, err := client.SayHello(context.TODO(), &connect.Request[proto.HelloRequest]{
		Msg: &proto.HelloRequest{
			Name: "Taha",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for stream.Receive() {
		msg := stream.Msg().Msg

		fmt.Println(msg)
	}
}
