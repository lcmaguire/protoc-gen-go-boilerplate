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

//go:embed templates/*
var embeddedTemplates embed.FS

const (
	// default dirPath
	defaultDir = "templates"

	unaryMethodSuffix        = "method.unary.go.tmpl"
	serverStreamMethodSuffix = "method.server.stream.go.tmpl"
	clientStreamMethodSuffix = "method.client.stream.go.tmpl"
	bidiStreamMethodSuffix   = "method.bidi.stream.go.tmpl"

	serviceSuffix = "service.go.tmpl"
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
		// if populated will override to be the provided directory.
		directory := defaultDir
		if directoryOverride != nil {
			directory = *directoryOverride
		}

		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}

			for _, service := range file.Services {
				serviceFileName := strings.ToLower(filepath.Join(service.GoName, "service.go"))
				sf := gen.NewGeneratedFile(serviceFileName, ".")
				sf.P("package " + file.GoPackageName)

				// imports
				ident := sf.QualifiedGoIdent(file.GoDescriptorIdent)
				sf.Import(file.GoImportPath)

				// generate service struct for go-grpc.
				// gets alias of file.GoDescriptorIdent
				pkgIdent := strings.Split(ident, ".")[0]

				methods := make([]Method, 0, len(service.Methods))

				for _, method := range service.Methods {
					fileName := strings.ToLower(filepath.Join(service.GoName, method.GoName+".go"))
					nf := gen.NewGeneratedFile(fileName, ".")
					nf.P("package " + file.GoPackageName)

					m := Method{
						MethodName:     method.GoName,
						MethodFullName: string(method.Desc.FullName()),
						ServiceName:    service.GoName,
						InputName:      messageImportPath(method.Input, nf),
						ResponseName:   messageImportPath(method.Output, nf),
						Ident:          pkgIdent,
						Method:         method,
						FileGoPkgName:  string(file.GoPackageName),
					}

					// get the appropriate suffix & the override template when applicable.
					methodSuffix := ""
					var overrideFile *string
					switch {
					case method.Desc.IsStreamingServer() && method.Desc.IsStreamingClient():
						methodSuffix = bidiStreamMethodSuffix
						overrideFile = bidiStreamMethodTemplate
					case method.Desc.IsStreamingServer():
						methodSuffix = serverStreamMethodSuffix
						overrideFile = serverStreamMethodTemplate
					case method.Desc.IsStreamingClient():
						methodSuffix = clientStreamMethodSuffix
						overrideFile = clientStreamMethodTemplate
					default:
						methodSuffix = unaryMethodSuffix
						overrideFile = customUnaryMethodTemplate
					}

					currentTemplate, err := loadTemplates(directory, methodSuffix, overrideFile)
					if err != nil {
						return err
					}

					buffy := bytes.NewBuffer([]byte{})
					if err := currentTemplate.Execute(buffy, m); err != nil {
						return err
					}

					nf.P(buffy.String())
					// will tidy the imports of the generated method file.
					err = tidyImports(gen, nf, fileName)
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

				serviceT, err := loadTemplates(directory, serviceSuffix, customServiceTemplate)
				if err != nil {
					return err
				}

				buffy := bytes.NewBuffer([]byte{})
				if err := serviceT.Execute(buffy, s); err != nil {
					return err
				}
				sf.P(buffy.String())

				// will tidy the imports of the generated service file.
				err = tidyImports(gen, sf, serviceFileName)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func loadTemplates(dir string, suffix string, override *string) (*template.Template, error) {
	if override != nil && len(*override) > 0 {
		return template.ParseFiles(*override)
	}

	return template.ParseFS(embeddedTemplates, filepath.Join(dir, suffix))
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
