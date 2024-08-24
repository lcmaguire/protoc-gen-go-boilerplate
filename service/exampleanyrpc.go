package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
	anypb "google.golang.org/protobuf/types/known/anypb"
)
import (
	"context"
)

// ExampleAnyRpc implements proto.ExampleAPI.ExampleAnyRpc.
func (s *Service) ExampleAnyRpc(ctx context.Context, in *temp.Example) (*anypb.Any, error) {
	return nil, nil
}
