package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

import (
	connect "connectrpc.com/connect"
	"context"
)

// ExampleRpc is a connect rpc implementation of proto.ExampleAPI.ExampleRpc.
func (s *Service) ExampleRpc(ctx context.Context, in *connect.Request[temp.Example]) (*connect.Response[temp.Example], error) {
	return nil, nil
}
