
// {{ .MethodName}} implements {{.MethodFullName}}.
func (s *Service) {{ .MethodName}} (in *{{.InputName}}, svr {{.Ident}}.{{.ServiceName}}_{{.MethodName}}Server) error {
    return nil
}