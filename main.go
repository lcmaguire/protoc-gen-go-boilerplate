package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {

	protogen.Options{
		//	ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {

		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}

			for _, service := range file.Services {

				fileImportPath := protogen.GoImportPath(".")

				sf := gen.NewGeneratedFile("service.go", fileImportPath)
				sf.P("package " + file.GoPackageName)

				// make sure this is imported for pkg.
				_ = sf.QualifiedGoIdent(file.GoDescriptorIdent)
				// generate go-grpc server.
				sName := string(file.GoPackageName) + "." + unimplementedServerString(service.GoName)

				s := Service{
					UnimplementedServiceName: sName,
				}
				sf.P(ExecuteTemplate(ServiceTemplate, s))

				//nf := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+".go", fileImportPath)
				//nf.P("package " + file.GoPackageName)

				//i := protogen.GoIdent{GoName: "Context", GoImportPath: protogen.GoImportPath("context")}
				//nf.QualifiedGoIdent(i)

				for _, method := range service.Methods {
					// snake case filename
					nf := gen.NewGeneratedFile(strings.ToLower(method.GoName)+".go", fileImportPath)
					nf.P("package " + file.GoPackageName)

					i := protogen.GoIdent{GoName: "Context", GoImportPath: protogen.GoImportPath("context")}
					nf.QualifiedGoIdent(i)
					m := Method{
						MethodName:   method.GoName,
						InputName:    messageImportPath(method.Input, nf),
						ResponseName: messageImportPath(method.Output, nf),
					}

					str := ExecuteTemplate(MethodTemplate, m)
					nf.P(str)
				}

			}
		}

		return nil
	})
}

// assumption this is always going to be in a different package.
func messageImportPath(in *protogen.Message, f *protogen.GeneratedFile) string {
	return f.QualifiedGoIdent(in.GoIdent)
}

func unimplementedServerString(serviceName string) string {
	return fmt.Sprintf("Unimplemented%sServer", serviceName)
}

// Service data regarding the service to be implemented.
type Service struct {
	UnimplementedServiceName string
}

// ServiceTemplate template represents the struct being generated to implement service.
const ServiceTemplate = `
type Service struct {
	{{.UnimplementedServiceName}}
}
`

// Method contains all info for method generation.
type Method struct {
	// MethodName is the name of the RPC being implemented.
	MethodName string

	// InputName import path and type name e.g foo.Bar.
	InputName string
	// InputName import path and type name e.g foo.Bar.
	ResponseName string
}

/*
	perhaps via config
	- functions  ()
*/

const MethodTemplate = `
	// {{ .MethodName}} ...
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

	func (s *Service) preDatabaseDownstreams{{ .MethodName}}(ctx context.Context, in any) (any, error){
		return nil, nil
	}

	func (s *Service) databaseOp{{ .MethodName}}(ctx context.Context, downstreamResponse any, internalType any) (any, error){
		return nil, nil
	}
	
	func (s *Service) postDatabaseDownstreams{{ .MethodName}}(ctx context.Context, in any) (any, error){
		return nil, nil
	}

	func validate{{ .MethodName}}Input(ctx context.Context, in *{{ .InputName}}) error {
		return nil
	}

	func map{{ .MethodName}}InputToInternal(ctx context.Context, in *{{ .InputName}}) (any, error) {
		return nil, nil
	}

	func prepare{{ .MethodName}}Response(ctx context.Context, downstreamResponse any, internalType any, databaseType any, postDbDownstreamResponse any) (*{{ .ResponseName}}, error) {
		return nil, nil
	}

`

// ExecuteTemplate something to implement templates.
func ExecuteTemplate(tplate string, data any) string {
	// todo read more about template library, see if it may be better to have one Template struct and re use it.
	templ, err := template.New("").Parse(tplate)
	if err != nil {
		panic(err)
	}

	buffy := bytes.NewBuffer([]byte{})
	if err := templ.Execute(buffy, data); err != nil {
		panic(err)
	}
	return buffy.String()
}
