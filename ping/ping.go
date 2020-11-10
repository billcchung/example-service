package Ping

import (
	"context"
	pb "github.com/billcchung/example-service/proto"
	"go.opencensus.io/trace"
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Server struct{}

func (s Server) Get(ctx context.Context, req *pb.PingRequest) (res *pb.PingResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "Get")
	defer span.End()
	res = &pb.PingResponse{
		Message_ID:  req.Message_ID,
		MessageBody: req.MessageBody,
		Timestamp:   uint64(time.Now().UnixNano() / int64(time.Millisecond)),
	}
	return
}

func (s Server) GetAfter(ctx context.Context, req *pb.PingRequestWithSleep) (res *pb.PingResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAfter")
	defer span.End()
	time.Sleep(time.Duration(req.Sleep) * time.Second)
	return s.Get(ctx, &pb.PingRequest{Message_ID: req.Message_ID, MessageBody: req.MessageBody})
}

func (s Server) GetRandom(ctx context.Context, req *pb.PingRequest) (res *pb.PingResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRandom")
	defer span.End()
	var garbage []string
	for i := 0; i <= 1000000; i++ {
		garbage = append(garbage, string(letterRunes[rand.Intn(len(letterRunes))]))
	}
	return s.Get(ctx, &pb.PingRequest{Message_ID: req.Message_ID, MessageBody: string(letterRunes[rand.Intn(len(letterRunes))])})
}
