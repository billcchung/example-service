package Ping

import (
	"context"
	"fmt"
	pb "github.com/billcchung/example-service/proto"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// call two functions in parallel
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return genGarbage(ctx)
	})

	var char string
	eg.Go(func() error {
		r, err := getRune(ctx, letterRunes)
		if err != nil {
			return err
		}
		char = string(r)
		return nil
	})

	// wait for the functions to finish
	if err := eg.Wait(); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("GetRandom err: %s", err))
	}

	return s.Get(ctx, &pb.PingRequest{Message_ID: req.Message_ID, MessageBody: char})
}

func getRune(ctx context.Context, runes []rune) (rune, error) {
	ctx, span := trace.StartSpan(ctx, "getRune")
	defer span.End()
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	return runes[rand.Intn(len(runes))], nil
}

func genGarbage(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "genGarbage")
	defer span.End()
	var garbage []string
	for i := 0; i <= rand.Intn(1000000)+100000; i++ {
		garbage = append(garbage, string(letterRunes[rand.Intn(len(letterRunes))]))
	}
	return nil
}
