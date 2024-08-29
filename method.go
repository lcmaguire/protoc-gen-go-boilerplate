package main

import (
	"embed"

	"google.golang.org/protobuf/compiler/protogen"
)

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
	// Method *protogen.Method.
	Method *protogen.Method
}

func (m Method) RunTemplate(defaultTemplate embed.FS, defaultPath string, filePath ...string) (string, error) {
	return generateTemplateData(m, defaultTemplate, defaultPath, filePath...)
}
