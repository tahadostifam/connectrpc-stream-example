package main

import (
	helloer "connectrpc-stream-example/gen/proto"
	"connectrpc-stream-example/gen/proto/protoconnect"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
)

type handler struct{}

func NewMyServiceHandler() protoconnect.MyProtoHandler {
	return &handler{}
}

// GiveInfo implements protoconnect.MyProtoHandler.
func (h *handler) GiveInfo(context.Context, *connect.Request[helloer.GiveInfoRequest]) (*connect.Response[helloer.GiveInfoResponse], error) {
	res := connect.NewResponse(&helloer.GiveInfoResponse{
		Data: "Hello from server!",
	})

	return res, nil
}

func (h *handler) SayHello(ctx context.Context, req *connect.Request[helloer.HelloRequest], stream *connect.ServerStream[helloer.HelloReply]) error {
	name := req.Msg.Name

	fmt.Println("Called")
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		stream.Send(&helloer.HelloReply{
			Msg: fmt.Sprintf("Hello %s", name),
		})
	}

	return nil
}

const port = 3000

func main() {
	api := http.NewServeMux()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	myServiceHandler := NewMyServiceHandler()
	path, handler := protoconnect.NewMyProtoHandler(myServiceHandler)
	api.Handle(path, handler)

	log.Println("Server listening at", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), api)
	if err != nil {
		log.Fatalln(err)
	}
}
