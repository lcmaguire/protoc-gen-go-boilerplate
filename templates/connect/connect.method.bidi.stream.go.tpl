import (
	connect "connectrpc.com/connect"
    "context"
)

// {{.MethodName}} is a connect rpc implementation of {{.MethodFullName}}.
func (s *Service) {{.MethodName}}(ctx context.Context, in *connect.BidiStream[{{.InputName}}, {{.ResponseName}}]) error {
	return nil
}