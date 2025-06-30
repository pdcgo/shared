package main

import (
	"context"
	"log"
	"net"

	"github.com/pdcgo/shared/interfaces/withdrawal_iface"
	"google.golang.org/grpc"
)

type MockWdServer struct {
	withdrawal_iface.UnimplementedWithdrawalServiceServer
}

// SubmitWithdrawal implements withdrawal_iface.WithdrawalServiceServer.
func (m *MockWdServer) SubmitWithdrawal(context.Context, *withdrawal_iface.SubmitWdRequest) (*withdrawal_iface.CommonResponse, error) {
	data := withdrawal_iface.CommonResponse{
		Message: "success",
	}

	return &data, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	withdrawal_iface.RegisterWithdrawalServiceServer(s, &MockWdServer{})

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
