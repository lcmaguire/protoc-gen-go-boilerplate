package temp

import (
	"context"

	"github.com/lcmaguire/protoc-gen-go-boilerplate/gen/temp"
)

// ExampleRpc implements proto.ExampleAPI.ExampleRpc.
func (s *Service) ExampleRpc(ctx context.Context, in *temp.Example) (*temp.Example, error) {
	// validate request
	err := validateExampleRpcInput(ctx, in)
	if err != nil {
		return nil, err
	}

	// map to internal type
	internalType, err := mapExampleRpcInputToInternal(ctx, in)
	if err != nil {
		return nil, err
	}

	// perform any dowsntream requests prior to database interaction.
	downstreamResponse, err := s.preDatabaseDownstreamsExampleRpc(ctx, internalType)
	if err != nil {
		return nil, err
	}

	// perform database operation
	databaseResponse, err := s.databaseOpExampleRpc(ctx, downstreamResponse, internalType)
	if err != nil {
		return nil, err
	}

	// perform any dowsntream requests post database interaction.
	postDbDownstreamResponse, err := s.postDatabaseDownstreamsExampleRpc(ctx, databaseResponse)
	if err != nil {
		return nil, err
	}

	// prepare response
	return prepareExampleRpcResponse(ctx, internalType, downstreamResponse, databaseResponse, postDbDownstreamResponse)
}

func (s *Service) preDatabaseDownstreamsExampleRpc(ctx context.Context, in any) (any, error) {
	return nil, nil
}

func (s *Service) databaseOpExampleRpc(ctx context.Context, downstreamResponse any, internalType any) (any, error) {
	return nil, nil
}

func (s *Service) postDatabaseDownstreamsExampleRpc(ctx context.Context, in any) (any, error) {
	return nil, nil
}

func validateExampleRpcInput(ctx context.Context, in *temp.Example) error {
	return nil
}

func mapExampleRpcInputToInternal(ctx context.Context, in *temp.Example) (any, error) {
	return nil, nil
}

func prepareExampleRpcResponse(ctx context.Context, downstreamResponse any, internalType any, databaseType any, postDbDownstreamResponse any) (*temp.Example, error) {
	return nil, nil
}
