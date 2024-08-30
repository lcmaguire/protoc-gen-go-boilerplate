package temp

import (
	"context"

	temp "github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// ExampleAnyRpc implements proto.ExampleAPI.ExampleAnyRpc.
func (s *Service) ExampleAnyRpc(ctx context.Context, in *temp.Example) (*anypb.Any, error) {
	// validate request
	err := validateExampleAnyRpcInput(ctx, in)
	if err != nil {
		return nil, err
	}

	// map to internal type
	internalType, err := mapExampleAnyRpcInputToInternal(ctx, in)
	if err != nil {
		return nil, err
	}

	// perform any dowsntream requests prior to database interaction.
	downstreamResponse, err := s.preDatabaseDownstreamsExampleAnyRpc(ctx, internalType)
	if err != nil {
		return nil, err
	}

	// perform database operation
	databaseResponse, err := s.databaseOpExampleAnyRpc(ctx, downstreamResponse, internalType)
	if err != nil {
		return nil, err
	}

	// perform any dowsntream requests post database interaction.
	postDbDownstreamResponse, err := s.postDatabaseDownstreamsExampleAnyRpc(ctx, databaseResponse)
	if err != nil {
		return nil, err
	}

	// prepare response
	return prepareExampleAnyRpcResponse(ctx, internalType, downstreamResponse, databaseResponse, postDbDownstreamResponse)
}

func (s *Service) preDatabaseDownstreamsExampleAnyRpc(ctx context.Context, in any) (any, error) {
	return nil, nil
}

func (s *Service) databaseOpExampleAnyRpc(ctx context.Context, downstreamResponse any, internalType any) (any, error) {
	return nil, nil
}

func (s *Service) postDatabaseDownstreamsExampleAnyRpc(ctx context.Context, in any) (any, error) {
	return nil, nil
}

func validateExampleAnyRpcInput(ctx context.Context, in *temp.Example) error {
	return nil
}

func mapExampleAnyRpcInputToInternal(ctx context.Context, in *temp.Example) (any, error) {
	return nil, nil
}

func prepareExampleAnyRpcResponse(ctx context.Context, downstreamResponse any, internalType any, databaseType any, postDbDownstreamResponse any) (*anypb.Any, error) {
	return nil, nil
}
