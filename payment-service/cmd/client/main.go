package main

import (
	"context"
	"log"
	"time"

	pb "go-grpc-payment/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Use WithTransportCredentials(insecure) because TLS/SSL certificates not setup
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.PaymentRequest{
		OrderId:        101,
		Amount:         99.50,
		Currency:       "USD",
		CredCardNumber: "1234-5678-9012-3456",
	}

	res, err := client.ProcessPayment(ctx, req)
	if err != nil {
		log.Fatalf("could not process payment: %v", err)
	}
	log.Printf("Response: Success=%v, TxID=%s, Msg=%s", res.Success, res.TransactionId, res.Message)
}
