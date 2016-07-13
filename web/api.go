package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/captncraig/erroneous"
)

type server struct {
	logger *erroneous.ErrorLogger
}

// GetMux returns an http.ServeMux that has all of the built-in routes for erroneous. The main error list will be at /{prefix} and adittional routes will be at /{prefix}/api/...adittional
//
// If no prefix is provided, the default value of /errors will be used
func GetMux(logger *erroneous.ErrorLogger, prefix string) *http.ServeMux {
	if prefix == "" {
		prefix = "/errors"
	}
	prefix = strings.TrimSuffix(prefix, "/")
	server := &server{logger}
	mux := http.NewServeMux()
	mux.HandleFunc(prefix, server.list)
	return mux
}

func (s *server) list(w http.ResponseWriter, r *http.Request) {
	errors, err := s.logger.GetAllErrors("")
	fmt.Println(listTpl.Execute(w, map[string]interface{}{"Count": len(errors), "Errors": errors, "Error": err}))
}

const listTemplate = `
{{define "title"}}{{.Count}} Errors{{end}}
{{define "content"}}
    {{if eq .Count 0}} No Errors Yet!
    {{else}}
        <h2 id="errorcount">ApplicationName - {{.Count}} Errors; last 2 minutes ago</h2>
        <table id="ErrorLog" class="alt-rows">
            <thead>
                <tr>
                    <th class="type-col">&nbsp;</th>
                    <th>Error</th>
                    <th>Url</th>
                    <th>Remote IP</th>
                    <th>Time</th>
                    <th>Site</th>
                    <th>Server</th>
                </tr>
            </thead>
            <tbody>
            {{range .Errors}}
                <tr>
                    <td></td>
                </tr>
            {{end}}
            </tbody>
        </table>
    {{end}}
{{end}}
`
const singleTemplate = ``
const baseTemplate = `
<html>
<head>
<title>Erroneous - {{block "title" .}}{{end}}</title>
</head>
<body>{{block "content" .}}No content?{{end}}</body>
</html>
`

var listTpl = template.Must(template.New("base").Parse(baseTemplate))
var singleTpl = template.Must(template.New("base").Parse(baseTemplate))

func init() {
	template.Must(listTpl.New("list").Parse(listTemplate))
	template.Must(singleTpl.New("single").Parse(singleTemplate))
}
