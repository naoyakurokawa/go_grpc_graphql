package controller

import (
	pb "github.com/naoyakurokawa/go_grpc_graphql_proto/pb"
	"google.golang.org/grpc"
)

// RegisterTaskService wires the Task service into the provided gRPC server.
func RegisterTaskService(grpcServer *grpc.Server, service pb.TaskServiceServer) {
	pb.RegisterTaskServiceServer(grpcServer, service)
}
