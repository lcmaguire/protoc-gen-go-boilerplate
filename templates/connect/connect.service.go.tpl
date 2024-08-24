
import (
 {{.ConnectGoImportPath}}
)

// Service connect implementation of {{.ServerFullName}}.
type Service struct {
{{.ServiceGoPkg}}connect.Unimplemented{{.ServiceName}}Handler
}
