package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"

	"context"

	connect "connectrpc.com/connect"
)

// ExampleServerStream implements ExampleServerStream
func (s *Service) ExampleServerStream(ctx context.Context, in *connect.Request[temp.Example]) (*connect.ServerStream[temp.Example], error) {
	return nil, nil
}
