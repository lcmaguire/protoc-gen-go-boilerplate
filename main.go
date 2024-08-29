package main

import (
	"bytes"
	"embed"
	"flag"
	"go/format"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
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

//go:embed templates/method.bidi.stream.go.tpl
var defaultBidiStreamingMethod embed.FS

//go:embed templates/connect
var connectTemplates embed.FS

//go:embed templates/service.go.tpl
var defaultServiceTemplate embed.FS

const (
	defaultUnaryMethodFilePath           = "templates/method.unary.go.tpl"
	defaultClientStreamingMethodFilePath = "templates/method.client.stream.go.tpl"
	defaultServerStreamingMethodFilePath = "templates/method.server.stream.go.tpl"
	defaultBidiStreamingMethodFilePath   = "templates/method.bidi.stream.go.tpl"
	defaultServiceFilePath               = "templates/service.go.tpl"

	defaultConnectUnaryMethodFilePath           = "templates/connect/connect.method.unary.go.tpl"
	defaultConnectClientStreamingMethodFilePath = "templates/connect/connect.method.client.stream.go.tpl"
	defaultConnectServerStreamingMethodFilePath = "templates/connect/connect.stream.go.tpl"
	defaultConnectBidiStreamingMethodFilePath   = "templates//connect/connect.stream.go.tpl"
	defaultConnectServiceFilePath               = "templates/connect/connect.service.go.tpl"
)

func main() {
	var flags flag.FlagSet
	customUnaryMethodTemplate := flags.String("unaryMethodTemplate", "", "custom method template")
	clientStreamMethodTemplate := flags.String("clientStreamMethodTemplate", "", "custom method template")
	serverStreamMethodTemplate := flags.String("serverStreamMethodTemplate", "", "custom method template")
	bidiStreamMethodTemplate := flags.String("bidiStreamMethodTemplate", "", "custom method template")
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
				sf := gen.NewGeneratedFile("service.go", ".")
				sf.P("package " + file.GoPackageName)

				// imports
				ident := sf.QualifiedGoIdent(file.GoDescriptorIdent)
				sf.Import(file.GoImportPath)

				// generate service struct for go-grpc.
				pkgIdent := strings.Split(ident, ".")[0]

				methods := make([]Method, 0, len(service.Methods))

				for _, method := range service.Methods {
					nf := gen.NewGeneratedFile(strings.ToLower(method.GoName)+".go", ".")
					nf.P("package " + file.GoPackageName)

					m := Method{
						MethodName:     method.GoName,
						MethodFullName: string(method.Desc.FullName()),
						ServiceName:    service.GoName,
						InputName:      messageImportPath(method.Input, nf),
						ResponseName:   messageImportPath(method.Output, nf),
						Ident:          pkgIdent,
						Method:         method,
					}

					// Will handle selecting appropriate template for method(s)
					// based upon if they are unary or a streaming rpc.
					methodTemplate := defaultUnaryMethodTemplate
					methodFile := defaultUnaryMethodFilePath
					override := ""
					if customUnaryMethodTemplate != nil {
						override = *customUnaryMethodTemplate
					}

					switch {
					case method.Desc.IsStreamingServer() && method.Desc.IsStreamingClient():
						methodFile = defaultBidiStreamingMethodFilePath
						methodTemplate = defaultBidiStreamingMethod
						if bidiStreamMethodTemplate != nil {
							override = *bidiStreamMethodTemplate
						}
					case method.Desc.IsStreamingServer():
						methodFile = defaultServerStreamingMethodFilePath
						methodTemplate = defaultServerStreamingMethod
						if clientStreamMethodTemplate != nil {
							override = *serverStreamMethodTemplate
						}
					case method.Desc.IsStreamingClient():
						methodFile = defaultClientStreamingMethodFilePath
						methodTemplate = defaultClientStreamingMethod
						if serverStreamMethodTemplate != nil {
							override = *clientStreamMethodTemplate
						}
					}

					methodBites, err := m.RunTemplate(methodTemplate, methodFile, override)
					if err != nil {
						return err
					}
					nf.P(methodBites)
					// will tidy the imports of the generated method file.
					err = tidyImports(gen, nf, imports.LocalPrefix+strings.ToLower(method.GoName)+".go")
					if err != nil {
						return err
					}

					methods = append(methods, m)
				}

				s := Service{
					ServiceGoImportPath: file.GoDescriptorIdent.String(),
					ConnectGoImportPath: connectPath(file).String(),
					FileGoPkgName:       string(file.GoPackageName),
					ServiceGoPkg:        pkgIdent,
					ServiceName:         service.GoName,
					Ident:               pkgIdent,
					ServerFullName:      string(service.Desc.FullName()),
					Methods:             methods,
					Service:             service,
				}

				serviceBites, err := s.RunTemplate(*customServiceTemplate)
				if err != nil {
					return err
				}
				sf.P(serviceBites)

				// will tidy the imports of the generated service file.
				err = tidyImports(gen, sf, imports.LocalPrefix+"service.go")
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// tidyImports will format imports into one import group and will remove any unused imports.
//
// does this via skipping the provided generatedFile & recreating an identical file with the same content but formatted.
func tidyImports(gen *protogen.Plugin, generatedFile *protogen.GeneratedFile, fileName string) error {
	bites, err := generatedFile.Content()
	if err != nil {
		return err
	}

	bites, err = format.Source(bites)
	if err != nil {
		return err
	}

	bites, err = imports.Process(fileName, bites, nil) // opt nil will result in default behaviour.
	if err != nil {
		return err
	}
	generatedFile.Skip()

	// will recreate an identical file but with the imports in order.
	newFile := gen.NewGeneratedFile(fileName, ".")
	newFile.P(string(bites))
	return nil
}

// assumption this is always going to be in a different package.
func messageImportPath(in *protogen.Message, f *protogen.GeneratedFile) string {
	return f.QualifiedGoIdent(in.GoIdent)
}

func generateTemplateData(data any, embededFile embed.FS, defaultPath string, filePath ...string) (string, error) {
	var templ *template.Template
	var err error

	// if no override is provided will use default.
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
