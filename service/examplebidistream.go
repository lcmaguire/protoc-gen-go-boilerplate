package temp

import (
	"github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

// ExampleBidiStream implements proto.ExampleAPI.ExampleBidiStream.
func (s *Service) ExampleBidiStream(in *temp.Example, svr temp.ExampleAPI_ExampleBidiStreamServer) error {
	return nil
}
