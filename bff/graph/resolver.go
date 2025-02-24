package graph

import (
	pb "github.com/naoyakurokawa/go_grpc_graphql_proto/pb"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	GrpcClient pb.TaskServiceClient // gRPC クライアントを追加
}
