package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/naoyakurokawa/go_grpc_graphql/Infrastructure/store"
	"github.com/naoyakurokawa/go_grpc_graphql/controller"
	"github.com/naoyakurokawa/go_grpc_graphql/graph"
	"github.com/naoyakurokawa/go_grpc_graphql/graph/resolver"
	"github.com/naoyakurokawa/go_grpc_graphql/pkg/pb"
	"github.com/naoyakurokawa/go_grpc_graphql/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcAddress = "backend:50051"

func main() {
	// gRPC クライアントの接続
	conn, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// gRPC クライアントを作成
	taskClient := pb.NewTaskServiceClient(conn)
	categoryClient := pb.NewCategoryServiceClient(conn)

	todoRepo := store.NewTodoStore(taskClient)
	categoryRepo := store.NewCategoryStore(categoryClient)

	todoUsecase := usecase.NewTodoUsecase(todoRepo)
	todoController := controller.NewTodoController(todoUsecase)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryController := controller.NewCategoryController(categoryUsecase)

	e := echo.New()

	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	graphqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: &resolver.Resolver{
				TodoController:     todoController,
				CategoryController: categoryController,
			}},
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
