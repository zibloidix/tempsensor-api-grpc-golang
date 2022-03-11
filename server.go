package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/zibloidix/tempsensor-api-grpc-golang/tempsensorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) GetData(req *tempsensorpb.DataRequest, stream tempsensorpb.TempSensorService_GetDataServer) error {
	s := req.GetSession()
	f := req.GetFormat()
	var resf tempsensorpb.DataResponse_Format

	if f == tempsensorpb.DataRequest_F {
		resf = tempsensorpb.DataResponse_F
	} else {
		resf = tempsensorpb.DataResponse_C
	}

	for i := 0; i < 100; i++ {
		res := &tempsensorpb.DataResponse{
			Session: s,
			Temp:    1.0 + rand.Float32()*(100.0-1.0),
			Format:  resf,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Microsecond)
	}

	return nil
}

func main() {
	fmt.Println("TempService start - I am sensor")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fail to listen: %v", err)
	}
	s := grpc.NewServer()
	tempsensorpb.RegisterTempSensorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
