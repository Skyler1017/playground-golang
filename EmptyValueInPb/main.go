package main

import (
	"context"
	"git.code.oa.com/trpc-go/trpc-go"
	"log"
	pb "playground/EmptyValueInPb/foo"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, req *pb.HelloReq, rsp *pb.HelloRsp) error {
	age := req.GetAge()
	if age == nil {
		log.Println("age is nil")
		return nil
	}
	log.Printf("age: %+v\n", age)
	log.Println("value: ", int(age.Value))
	return nil
}

func main() {
	ts := trpc.NewServer()
	s := &Server{}
	pb.RegisterHelloService(ts, s)
	_ = ts.Serve()
}
