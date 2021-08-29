package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	gp "github.com/carloscontrerasruiz/grpc-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var requests = []*gp.LongGreetRequest{
	{
		Greeting: &gp.Greeting{
			FirstName: "Carlos",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "Lucy",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "Tom",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "CC",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "DD",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "Cinthya",
		},
	},
}

var requestsEvery = []*gp.GreetEveryoneRequest{
	{
		Greeting: &gp.Greeting{
			FirstName: "Carlos",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "Lucy",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "Tom",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "CC",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "DD",
		},
	},
	{
		Greeting: &gp.Greeting{
			FirstName: "Cinthya",
		},
	},
}

func main() {
	fmt.Println("Hello Im the client")

	opts := grpc.WithInsecure()

	tls := false

	if tls {
		certFile := "../ssl/ca.crt"

		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")

		if sslErr != nil {
			log.Fatalf("Error certificate client %v", sslErr)
		}

		opts = grpc.WithTransportCredentials(creds)
	}

	//opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Could not conect: %v", err)
	}

	defer cc.Close()

	c := gp.NewGreetServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)

	//doClientStreamin(c)

	//doByDiStreamin(c)

	doUnaryWithDeadline(c)

}
func doByDiStreamin(c gp.GreetServiceClient) {
	//create de client
	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while calling longgreet %v", err)
	}

	waitc := make(chan struct{})
	//we send a bunch of message
	go func() {
		for _, req := range requestsEvery {
			fmt.Println("sending message " + req.Greeting.FirstName)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//we receive a lot of messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error greeteveryone greet %v", err)
				break
			}

			fmt.Println("receving " + res.Result)
		}

		close(waitc)
	}()
	//block until everyting is done
	<-waitc
}

func doClientStreamin(c gp.GreetServiceClient) {

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling longgreet %v", err)
	}

	for _, req := range requests {
		//we sen each message individual
		fmt.Println("Sending request")
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while response longgreet %v", err)
	}

	fmt.Println(response)

}

func doServerStreaming(c gp.GreetServiceClient) {

	req := &gp.GreetManyTimesRequest{
		Greeting: &gp.Greeting{
			FirstName: "Carlos",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling greetmanytimes %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while calling greetmanytimes %v", err)
		}
		log.Printf("Response from gretmanytimes %v", msg.Result)
	}
}
func doUnary(c gp.GreetServiceClient) {
	req := &gp.GreetRequest{
		Greeting: &gp.Greeting{
			FirstName: "Carlos111",
			LastName:  "CC",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling greet %v", err)
	}

	log.Printf("Response from greeting %v", res.Result)
}

func doUnaryWithDeadline(c gp.GreetServiceClient) {
	req := &gp.GreetWithDeadlineRequest{
		Greeting: &gp.Greeting{
			FirstName: "Carlos111",
			LastName:  "CC",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Deadline was exceded")
			} else {
				fmt.Printf("Error during calling greetwithdeadline %v", err)
			}
		} else {
			log.Fatalf("Error while calling greet %v", err)
		}
		return
	}

	log.Printf("Response from greeting with dead line %v", res.Result)
}
