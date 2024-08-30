package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

// Method contains all info for method generation.
type Method struct {
	// MethodName is the name of the RPC being implemented.
	MethodName string
	// MethodFullName full method name.
	//
	// example.service.MethodName
	MethodFullName string
	// ServiceName is the name of the service to which the method belongs.
	ServiceName string
	// Ident the file pkg name.
	Ident string
	// InputName import path and type name e.g foo.Bar.
	InputName string
	// ResponseName import path and type name for the rpc response e.g foo.Bar.
	ResponseName string
	// Method *protogen.Method.
	Method *protogen.Method
}

// RunTemplate will execute the template.
func (m Method) RunTemplate(defaultPath string, filePath ...string) (string, error) {
	return generateTemplateData(m, defaultPath, filePath...)
}
