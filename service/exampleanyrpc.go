package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

import (
	connect "connectrpc.com/connect"
	"context"
)

// ExampleAnyRpc is a connect rpc implementation of proto.ExampleAPI.ExampleAnyRpc.
func (s *Service) ExampleAnyRpc(ctx context.Context, in *connect.Request[temp.Example]) (*connect.Response[anypb.Any], error) {
	return nil, nil
}
