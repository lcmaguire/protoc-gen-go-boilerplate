package temp

import (
	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

// ExampleClientStream implements proto.ExampleAPI.ExampleClientStream.
func (s *Service) ExampleClientStream(in temp.ExampleAPI_ExampleClientStreamServer) error {
	return nil
}
