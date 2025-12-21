package resolver

import (
	"github.com/naoyakurokawa/go_grpc_graphql/controller"
	"github.com/naoyakurokawa/go_grpc_graphql/middleware/session"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoController     *controller.TodoController
	CategoryController *controller.CategoryController
	AuthController     *controller.AuthController
	SessionManager     *session.Manager
}
