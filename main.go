package main

import (
	"context"
	"log"
	"net"
	"sync"

	// import the generated protobuf code
	pb "github.com/septilsnna/ms-example/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository - Dummy repository, this simulates the use of a datastore of some kind.
// We'll replace this with a real implementation later on.
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

// Service should implement all of the methods to statisfy the service we defined in our
// protobuf definition. You can check the interface in the generated code itself for
// the exact method signatures etc to give you a better idea.
type service struct {
	repo repository
}

// CreateConsignment - we created just one method on our service, which is a create methode,
// which takes a context and a request as an argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the 'Response' message we created in our protobuf definition
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	// Set up our gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with gRPC serveer, this will tie our implementation into the
	// auto-generated interface code for our protobuf definition.
	pb.RegisterShippingServicesServer(s, &service{repo})

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("faied to serve: %v", err)
	}
}
