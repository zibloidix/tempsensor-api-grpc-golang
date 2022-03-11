package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/zibloidix/tempsensor-api-grpc-golang/tempsensorpb"
	"google.golang.org/grpc"
)

func main() {
	session := flag.String("session", "304eaac5-3b3b-4adb-be53-0f295eb8f5a5", "Session UUID for Client")
	format := flag.String("format", "C", "Temperature format: C - Celsius, F - Fahrenheit")
	flag.Parse()
	fmt.Println("Client start: " + *session)

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()
	c := tempsensorpb.NewTempSensorServiceClient(cc)

	f := tempsensorpb.DataRequest_C
	if *format == "F" {
		f = tempsensorpb.DataRequest_F
	}
	req := &tempsensorpb.DataRequest{
		Session: *session,
		Format:  f,
	}

	resStream, err := c.GetData(context.Background(), req)
	if err != nil {
		log.Fatalf("Error call GetData(): %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// Service stream end
			break
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}
		msgf := msg.GetFormat()
		msgs := msg.GetSession()
		msgt := msg.GetTemp()
		log.Printf("Sensor return data. Session: %v, temp: %v, format: %v", msgs, msgt, msgf)
	}
}
