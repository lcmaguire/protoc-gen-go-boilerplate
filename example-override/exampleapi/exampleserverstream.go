package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

// ExampleServerStream implements proto.ExampleAPI.ExampleServerStream.
func (s *Service) ExampleServerStream(in *temp.Example, svr temp.ExampleAPI_ExampleServerStreamServer) error {
	return nil
}
