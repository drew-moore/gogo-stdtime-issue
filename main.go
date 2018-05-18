package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"net"
	"time"
)

func main() {
	runServer()
	runClient()
}

const (
	serverHost = "localhost"
	serverPort = 4321
)

/********************************  client  **************************************/
func runClient() {
	conn, err := grpc.DialContext(context.TODO(), fmt.Sprintf("%s:%d", serverHost, serverPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	client := NewThingServiceClient(conn)

	thing, err := client.GetThing(context.TODO(), &empty.Empty{})

	if err != nil {
		panic(err)
	}

	spew.Dump(thing)

}

/******************************** grpc service impl *******************************/
type thingServer struct{}

func (t *thingServer) GetThing(context.Context, *empty.Empty) (*Thing, error) {
	now := time.Now()
	return &Thing{
		Created: &now,
	}, nil
}

/********************************  server  **************************************/

func runServer() {
	address, _ := net.ResolveTCPAddr(`tcp`, fmt.Sprintf("%s:%d", serverHost, serverPort))
	listener, _ := net.ListenTCP(`tcp4`, address)

	serviceImpl := thingServer{}

	server := grpc.NewServer()

	RegisterThingServiceServer(server, &serviceImpl)

	go func() {
		err := server.Serve(listener)

		// returns when lis.Accept fails with fatal errors.
		if err != nil {
			panic(err)
		}
	}()

}
