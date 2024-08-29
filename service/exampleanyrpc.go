package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
	anypb "google.golang.org/protobuf/types/known/anypb"

	"context"

	connect "connectrpc.com/connect"
)

// ExampleAnyRpc is a connect rpc implementation of proto.ExampleAPI.ExampleAnyRpc.
func (s *Service) ExampleAnyRpc(ctx context.Context, in *connect.Request[temp.Example]) (*connect.Response[anypb.Any], error) {
	return nil, nil
}
