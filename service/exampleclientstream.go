package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"

	"context"

	connect "connectrpc.com/connect"
)

// ExampleClientStream implements ExampleClientStream
func (s *Service) ExampleClientStream(ctx context.Context, in *connect.ClientStream[temp.Example]) (*connect.Response[temp.Example], error) {
	return nil, nil
}
