package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

import (
	connect "connectrpc.com/connect"
	"context"
)

// ExampleClientStream implements ExampleClientStream
func (s *Service) ExampleClientStream(ctx context.Context, in *connect.Request[temp.Example]) (*connect.ServerStream[temp.Example], error) {
	return nil, nil
}
