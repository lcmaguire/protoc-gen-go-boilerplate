package main

import (
	"bytes"
	"embed"
	"flag"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// todo will need custom register of embeded files for types.

//go:embed templates/method.unary.go.tpl
var defaultUnaryMethodTemplate embed.FS

//go:embed templates/method.client.stream.go.tpl
var defaultClientStreamingMethod embed.FS

//go:embed templates/method.server.stream.go.tpl
var defaultServerStreamingMethod embed.FS

//go:embed templates/method.server.stream.go.tpl
var defaultBidiStreamingMethod embed.FS

//go:embed service.go.tpl
var defaultServiceTemplate embed.FS

const (
	defaultUnaryMethodFilePath           = "templates/method.unary.go.tpl"
	defaultClientStreamingMethodFilePath = "templates/method.client.stream.go.tpl"
	defaultServerStreamingMethodFilePath = "templates/method.server.stream.go.tpl"
	defaultBidiStreamingMethodFilePath   = "templates/method.server.stream.go.tpl"

	defaultServiceFilePath = "service.go.tpl"
)

func main() {
	var flags flag.FlagSet
	// todo have multiple flags for this.
	customUnaryMethodTemplate := flags.String("unaryMethodTemplate", "", "custom method template")
	clientStreamMethodTemplate := flags.String("clientStreamMethodTemplate", "", "custom method template")
	serverStreamMethodTemplate := flags.String("serverStreamMethodTemplate", "", "custom method template")
	bidiStreamMethodTemplate := flags.String("bidiStreamMethodTemplate", "", "custom method template")

	//
	customServiceTemplate := flags.String("serviceTemplate", "", "custom service template")

	protogen.Options{
		ParamFunc: flags.Set,
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
				// todo seek a way to ommit this for connect / non standard gRPC services.
				ident := sf.QualifiedGoIdent(file.GoDescriptorIdent)
				sf.Import(file.GoImportPath) // import regardless

				// generate service struct for go-grpc.
				pkgIdent := strings.Split(ident, ".")[0]

				s := Service{
					ServiceGoImportPath: file.GoDescriptorIdent.String(),
					ConnectGoImportPath: connectPath(file).String(),
					FileGoPkgName:       string(file.GoPackageName),
					ServiceGoPkg:        pkgIdent,
					ServiceName:         service.GoName,
					Ident:               pkgIdent,
					ServerFullName:      string(service.Desc.FullName()),
				}

				serviceBites, err := s.RunTemplate(*customServiceTemplate)
				if err != nil {
					return err
				}
				sf.P(serviceBites)

				for _, method := range service.Methods {
					// snake case filename
					nf := gen.NewGeneratedFile(strings.ToLower(method.GoName)+".go", fileImportPath)
					nf.P("package " + file.GoPackageName)

					//i := protogen.GoIdent{GoName: "Context", GoImportPath: protogen.GoImportPath("context")}
					//nf.QualifiedGoIdent(i)

					m := Method{
						MethodName:     method.GoName,
						MethodFullName: string(method.Desc.FullName()),
						ServiceName:    service.GoName,
						InputName:      messageImportPath(method.Input, nf),
						ResponseName:   messageImportPath(method.Output, nf),
						Ident:          pkgIdent,
					}

					// Will handle selecting appropriate template for method(s)
					// based upon if they are unary or a streaming rpc.
					methodTemplate := defaultUnaryMethodTemplate
					methodFile := defaultUnaryMethodFilePath
					overide := ""
					if customUnaryMethodTemplate != nil {
						overide = *customUnaryMethodTemplate
					}

					switch {
					case method.Desc.IsStreamingServer() && method.Desc.IsStreamingClient():
						methodFile = defaultBidiStreamingMethodFilePath
						methodTemplate = defaultBidiStreamingMethod
						if bidiStreamMethodTemplate != nil {
							overide = *bidiStreamMethodTemplate
						}
					case method.Desc.IsStreamingServer():
						methodFile = defaultServerStreamingMethodFilePath
						methodTemplate = defaultServerStreamingMethod
						if clientStreamMethodTemplate != nil {
							overide = *serverStreamMethodTemplate
						}
					case method.Desc.IsStreamingClient():
						methodFile = defaultClientStreamingMethodFilePath
						methodTemplate = defaultClientStreamingMethod
						if serverStreamMethodTemplate != nil {
							overide = *clientStreamMethodTemplate
						}
					}

					methodBites, err := m.RunTemplate(methodTemplate, methodFile, overide)
					if err != nil {
						return err
					}
					nf.P(methodBites)
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

// Service data regarding the service to be implemented.
type Service struct {
	// ServiceGoImportPath used for services
	ServiceGoImportPath string
	// ConnectGoImportPath generated connect import path.
	ConnectGoImportPath string
	// FileGoPkgName ...
	FileGoPkgName string
	// ServiceGoPkg last dir in package.
	ServiceGoPkg string
	// ServiceName
	ServiceName string
	// Ident the file pkg name.
	Ident string
	// ServerFullName full service name e.g foo.bar.service.
	ServerFullName string
}

func (s Service) RunTemplate(filePath ...string) (string, error) {
	return generateTemplateData(s, defaultServiceTemplate, defaultServiceFilePath, filePath...)
}

// Method contains all info for method generation.
type Method struct {
	// MethodName is the name of the RPC being implemented.
	MethodName string
	// MethodFullName full method name.
	MethodFullName string
	// ServiceName
	ServiceName string
	// Ident the file pkg name.
	Ident string

	// InputName import path and type name e.g foo.Bar.
	InputName string
	// InputName import path and type name e.g foo.Bar.
	ResponseName string
}

func (m Method) RunTemplate(defaultTemplate embed.FS, defaultPath string, filePath ...string) (string, error) {
	return generateTemplateData(m, defaultTemplate, defaultPath, filePath...)
}

func generateTemplateData(data any, embededFile embed.FS, defaultPath string, filePath ...string) (string, error) {
	var templ *template.Template
	var err error

	// if no override is provided use default.
	if len(filePath) == 0 || filePath[0] == "" {
		templ, err = template.ParseFS(embededFile, defaultPath)
	} else {
		templ, err = template.ParseFiles(filePath...)
	}

	if err != nil {
		return "", err
	}
	buffy := bytes.NewBuffer([]byte{})
	if err := templ.Execute(buffy, data); err != nil {
		return "", err
	}
	return buffy.String(), nil
}
