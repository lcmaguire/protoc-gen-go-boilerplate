package temp

import (
	"context"

	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

// ExampleRpc implements proto.ExampleAPI.ExampleRpc.
func (s *Service) ExampleRpc(ctx context.Context, in *temp.Example) (*temp.Example, error) {
	return nil, nil
}
