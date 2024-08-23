package main

import (
	"path"

	"google.golang.org/protobuf/compiler/protogen"
)

// connectPath used to get the connect path for a service file.
func connectPath(file *protogen.File) protogen.GoImportPath {
	connectFileName := file.GoPackageName + "connect"
	importP := protogen.GoImportPath(path.Join(
		string(file.GoImportPath),
		string(connectFileName),
	))
	return importP
}
