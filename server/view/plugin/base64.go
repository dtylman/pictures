package plugin

import (
	"encoding/base64"
	"html/template"
)

// NoEscape returns a template.FuncMap
// * NOESCAPE prevents escaping variable
func Base64() template.FuncMap {
	f := make(template.FuncMap)

	f["BASE64"] = func(name string) string {
		return base64.StdEncoding.EncodeToString([]byte(name))
	}

	return f
}
