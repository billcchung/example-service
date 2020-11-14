package main

import (
	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"flag"
	"github.com/billcchung/example-service/ping"
	pb "github.com/billcchung/example-service/proto"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	var projectID string
	var service string
	var serviceVersion string
	var address string
	flag.StringVar(&projectID, "p", "", "The GCP Project ID")
	flag.StringVar(&service, "s", "example-service", "The service name")
	flag.StringVar(&serviceVersion, "v", "1", "The service version")
	flag.StringVar(&address, "a", ":8080", "The address to listen for connections")
	flag.Parse()

	// Setup Cloud Profiler
	err := profiler.Start(profiler.Config{
		Service:        service,
		ServiceVersion: serviceVersion,
		ProjectID:      projectID,
	})
	if err != nil {
		log.Fatalf("Unable to start GCP profiler, err: %s", err)
	}
	log.Println("configured cloud profiler")

	// Setup Cloud Trace
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	// use ProbabilitySampler for production
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	log.Println("configured cloud trace")

	s := grpc.NewServer()
	pb.RegisterPingServer(s, Ping.Server{})
	if err != nil {
		log.Fatalf("Failed to get new gRPC server, err: %s", err)
		return
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen for address '%s', err: %s", address, err)
	}
	log.Println("Starting grpc server with cloud profiler", "project", projectID, "service", service, "addr", address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc, err: %s", err)
	}

}
