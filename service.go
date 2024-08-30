package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

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
	// Methods the methods for the service.
	Methods []Method
	// Service the protogen Service.
	Service *protogen.Service
}

// RunTemplate will execute the template.
func (s Service) RunTemplate(defaultPath string, filePath ...string) (string, error) {
	return generateTemplateData(s, defaultPath, filePath...)
}
