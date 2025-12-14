package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "go-grpc-payment/proto/payment"

	"google.golang.org/grpc"
)

const PORT = ":50051"

type server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *server) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("Received payment request for Order ID: %d, Amount: %.2f", req.OrderId, req.Amount)

	// Simulation logic
	if req.Amount > 1000 {
		return &pb.PaymentResponse{
			Success:       false,
			Message:       "Transaction limit exceeded",
			TransactionId: "",
		}, nil
	}

	return &pb.PaymentResponse{
		Success:       true,
		Message:       "Payment processed successfully",
		TransactionId: fmt.Sprintf("TXN-%d", req.OrderId),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, &server{})

	log.Printf("gRPC Payment Server listening on %s", PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
