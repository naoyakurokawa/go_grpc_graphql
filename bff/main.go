package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/naoyakurokawa/go_grpc_graphql/graph"
	pb "github.com/naoyakurokawa/go_grpc_graphql_proto/pb"
	"google.golang.org/grpc"
)

const grpcAddress = "backend:50051"

func main() {
	// gRPC クライアントの接続
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// gRPC クライアントを作成
	grpcClient := pb.NewTaskServiceClient(conn)

	e := echo.New()

	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	graphqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: &graph.Resolver{GrpcClient: grpcClient}},
		),
	)
	playgroundHandler := playground.Handler("GraphQL", "/query")

	e.POST("/query", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	err = e.Start(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}
