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

//go:embed method.go.tpl
var defaultMethodTemplate embed.FS

//go:embed service.go.tpl
var defaultServiceTemplate embed.FS

const (
	defaultMethodFilePath  = "method.go.tpl"
	defaultServiceFilePath = "service.go.tpl"
)

func main() {
	// todo set up custom templates.
	var flags flag.FlagSet
	customMethodTemplate := flags.String("methodTemplate", "", "custom method template")
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
				ident := sf.QualifiedGoIdent(file.GoDescriptorIdent)

				// generate service struct for go-grpc.
				pkgIdent := strings.Split(ident, ".")[0]

				s := Service{
					ServiceGoImportPath: string(file.GoDescriptorIdent.String()),
					ServiceGoPkg:        pkgIdent,
					Ident:               pkgIdent,
					ServiceName:         service.GoName,
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

					i := protogen.GoIdent{GoName: "Context", GoImportPath: protogen.GoImportPath("context")}
					nf.QualifiedGoIdent(i)
					m := Method{
						MethodName:     method.GoName,
						InputName:      messageImportPath(method.Input, nf),
						ResponseName:   messageImportPath(method.Output, nf),
						MethodFullName: string(method.Desc.FullName()),
					}

					methodBites, err := m.RunTemplate(*customMethodTemplate)
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

	// ConnectGoImportPath ...
	ConnectGoImportPath string

	// FileGoPkgName ...
	FileGoPkgName string

	ServiceGoPkg string

	// ServiceName
	ServiceName string
	// Ident the pkg name.
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
	// InputName import path and type name e.g foo.Bar.
	InputName string
	// InputName import path and type name e.g foo.Bar.
	ResponseName string
	// MethodFullName full method name.
	MethodFullName string
}

func (m Method) RunTemplate(filePath ...string) (string, error) {
	return generateTemplateData(m, defaultMethodTemplate, defaultMethodFilePath, filePath...)
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
