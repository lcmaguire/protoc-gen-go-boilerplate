// {{ .MethodName}} implements {{.MethodFullName}}.
func (s *Service) {{ .MethodName}}(ctx context.Context, in *{{ .InputName}} ) (*{{ .ResponseName}} , error) {
   	// validate request
   	err := validate{{ .MethodName}}Input(ctx, in)
   	if err != nil {
   		return nil, err
   	}

   	// map to internal type
   	internalType, err := map{{ .MethodName}}InputToInternal(ctx, in)
   	if err != nil {
   		return nil, err
   	}

   	// perform any dowsntream requests prior to database interaction.
   	downstreamResponse, err := s.preDatabaseDownstreams{{ .MethodName}}(ctx, internalType)
   	if err != nil {
   		return nil, err
   	}

   	// perform database operation
   	databaseResponse, err := s.databaseOp{{ .MethodName}}(ctx, downstreamResponse, internalType)
   	if err != nil {
   		return nil, err
   	}

   	// perform any dowsntream requests post database interaction.
   	postDbDownstreamResponse, err := s.postDatabaseDownstreams{{ .MethodName}}(ctx, databaseResponse)
   	if err != nil {
   		return nil, err
   	}

   	// prepare response
   	return prepare{{ .MethodName}}Response(ctx, internalType, downstreamResponse, databaseResponse, postDbDownstreamResponse)
}

func (s *Service) preDatabaseDownstreams{{.MethodName}}(ctx context.Context, in any) (any, error) {
	return nil, nil
}

func (s *Service) databaseOp{{.MethodName}}(ctx context.Context, downstreamResponse any, internalType any) (any, error) {
	return nil, nil
}

func (s *Service) postDatabaseDownstreams{{.MethodName}}(ctx context.Context, in any) (any, error) {
	return nil, nil
}

func validate{{.MethodName}}Input(ctx context.Context, in *temp.Example) error {
	return nil
}

func map{{.MethodName}}InputToInternal(ctx context.Context, in *temp.Example) (any, error) {
	return nil, nil
}

func prepare{{.MethodName}}Response(ctx context.Context, downstreamResponse any, internalType any, databaseType any, postDbDownstreamResponse any) (*{{ .ResponseName}}, error) {
	return nil, nil
}
