package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)
import (
	connect "connectrpc.com/connect"
	"context"
)

// ExampleBidiStream is a connect rpc implementation of proto.ExampleAPI.ExampleBidiStream.
func (s *Service) ExampleBidiStream(ctx context.Context, in *connect.BidiStream[temp.Example, temp.Example]) error {
	return nil
}
