package main

import (
	"context"
	"log"

	micro "github.com/Allenxuxu/mMicro"
	example "github.com/Allenxuxu/mMicro/example/srv/proto/hello"
	"github.com/Allenxuxu/mMicro/metadata"
	"github.com/Allenxuxu/mMicro/server"
	"github.com/Allenxuxu/mMicro/transport/grpc"
)

type Example struct{}

// Call(context.Context, *Request, *Response) error
//	Stream(context.Context, *StreamingRequest, Example_StreamStream) error
//	PingPong(context.Context, Example_PingPongStream) error
func (e *Example) Call(ctx context.Context, req *example.Request, rsp *example.Response) error {
	md, _ := metadata.FromContext(ctx)
	log.Printf("Received Example.Call request with metadata: %v", md)
	rsp.Msg = server.DefaultOptions().Id + ": Hello " + req.Name
	return nil
}

func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Print("Executing streaming handler")

	log.Printf("Received Example.Stream request with count: %d", req.Count)
	for i := 0; i < int(req.Count); i++ {
		log.Printf("Responding: %d", i)

		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	if err := stream.Close(); err != nil {
		return err
	}
	return nil
}

func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req := &example.Ping{}
		if err := stream.RecvMsg(req); err != nil {
			return err
		}
		log.Printf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

type Sub struct{}

func (e *Sub) Handle(ctx context.Context, msg *example.Message) error {
	log.Print("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *example.Message) error {
	log.Print("Function Received message: ", msg.Say)
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.example"),
		micro.Transport(grpc.NewTransport()),
	)

	// optionally setup command line usage
	service.Init()

	// Register Subscribers
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			"topic.example",
			new(Sub),
		),
	); err != nil {
		log.Fatal(err)
	}

	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			"topic.example",
			Handler,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Register Handlers
	example.RegisterExampleHandler(service.Server(), new(Example))

	// Run server
	if err := service.Run(); err != nil {
		panic(err)
	}
}
