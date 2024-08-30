package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

// ExampleBidiStream implements proto.ExampleAPI.ExampleBidiStream.
func (s *Service) ExampleBidiStream(svr temp.ExampleAPI_ExampleBidiStreamServer) error {
	return nil
}
