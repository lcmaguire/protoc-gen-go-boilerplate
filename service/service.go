package temp

import (
	"github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp/tempconnect"
)

// Service connect implementation of proto.ExampleAPI.
type Service struct {
	tempconnect.UnimplementedExampleAPIHandler
}
