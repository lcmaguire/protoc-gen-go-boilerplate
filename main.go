package main

import (
	"bytes"
	"embed"
	"flag"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// todo will need custom register of embeded files for types.

//go:embed templates/method.unary.go.tmpl
var defaultUnaryMethodTemplate embed.FS

//go:embed templates/service.go.tmpl
var defaultServiceTemplate embed.FS

//go:embed templates/*
var grpcTemplates embed.FS

//go:embed templates/connect/*
var connectTemplates embed.FS

const (
	// default dirPath
	defaultDir = "templates"

	unaryMethodSuffix        = "method.unary.go.tmpl"
	serverStreamMethodSuffix = "method.server.stream.go.tmpl"
	clientStreamMethodSuffix = "method.client.stream.go.tmpl"
	bidiStreamMethodSuffix   = "method.bidi.stream.go.tmpl"

	defaultServiceFilePath = "templates/service.go.tmpl"
	serviceSuffix          = "service.go.tmpl"
)

func main() {
	var flags flag.FlagSet
	customUnaryMethodTemplate := flags.String("unaryMethodTemplate", "", "custom method template")
	clientStreamMethodTemplate := flags.String("clientStreamMethodTemplate", "", "custom method template")
	serverStreamMethodTemplate := flags.String("serverStreamMethodTemplate", "", "custom method template")
	bidiStreamMethodTemplate := flags.String("bidiStreamMethodTemplate", "", "custom method template")
	customServiceTemplate := flags.String("serviceTemplate", "", "custom service template")

	directoryOverride := flags.String("templateDirectory", defaultDir, "custom directory for templates")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		// will user `templates` as the default dir.
		// if populated will overide to be the provided directory.
		directory := defaultDir
		if directoryOverride != nil {
			directory = *directoryOverride
		}

		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}

			for _, service := range file.Services {
				// todo have it be prefixed by service.GoName/
				sf := gen.NewGeneratedFile("service.go", ".")
				sf.P("package " + file.GoPackageName)

				// imports
				ident := sf.QualifiedGoIdent(file.GoDescriptorIdent)
				sf.Import(file.GoImportPath)

				// generate service struct for go-grpc.
				// gets alias of file.GoDescriptorIdent
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

					methodFile := filepath.Join(directory, unaryMethodSuffix)
					override := ""
					if customUnaryMethodTemplate != nil {
						override = *customUnaryMethodTemplate
					}

					switch {
					case method.Desc.IsStreamingServer() && method.Desc.IsStreamingClient():
						methodFile = filepath.Join(directory, bidiStreamMethodSuffix)
						if bidiStreamMethodTemplate != nil {
							override = *bidiStreamMethodTemplate
						}
					case method.Desc.IsStreamingServer():
						methodFile = filepath.Join(directory, serverStreamMethodSuffix)
						if clientStreamMethodTemplate != nil {
							override = *serverStreamMethodTemplate
						}
					case method.Desc.IsStreamingClient():
						methodFile = filepath.Join(directory, clientStreamMethodSuffix)
						if serverStreamMethodTemplate != nil {
							override = *clientStreamMethodTemplate
						}
					}

					methodBites, err := m.RunTemplate(methodFile, override)
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
					ServiceName:         service.GoName,
					ServerFullName:      string(service.Desc.FullName()),
					Methods:             methods,
					Service:             service,
					Ident:               pkgIdent,
				}

				serviceFile := filepath.Join(directory, serviceSuffix)

				serviceBites, err := s.RunTemplate(serviceFile, *customServiceTemplate)
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

func generateTemplateData(data any, defaultPath string, filePath ...string) (string, error) {
	var templ *template.Template
	var err error

	// if no override is provided will use default.
	if len(filePath) == 0 || filePath[0] == "" {
		//templ, err = template.ParseFS(embededFile, defaultPath)
		templ, err = template.ParseFiles(defaultPath)
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
