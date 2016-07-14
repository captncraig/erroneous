package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/captncraig/erroneous"
	"github.com/twinj/uuid"
)

type server struct {
	logger *erroneous.ErrorLogger
	prefix string
}

// GetMux returns an http.ServeMux that has all of the built-in routes for erroneous. The main error list will be at /{prefix} and adittional routes will be at /{prefix}/api/...adittional
//
// If no prefix is provided, the default value of /errors will be used
func GetMux(logger *erroneous.ErrorLogger, prefix string) *http.ServeMux {
	if prefix == "" {
		prefix = "/errors"
	}
	prefix = strings.TrimSuffix(prefix, "/")
	server := &server{logger, prefix}
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/", server.single)

	return mux
}

func (s *server) list(w http.ResponseWriter, r *http.Request) {
	errors, err := s.logger.GetAllErrors("")
	count := 0
	var lastTime time.Time
	for _, e := range errors {
		if e.CreationDate.After(lastTime) {
			lastTime = e.CreationDate
		}
		count += e.Count
	}
	ago := time.Now().Sub(lastTime)
	fmt.Println(listTpl.Execute(w, map[string]interface{}{"Count": count, "Errors": errors, "Msg": err, "Prefix": s.prefix, "Ago": ago.String()}))
}

func (s *server) single(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path == s.prefix+"/" { //can't distinguish between /errors and /errors/{id} with built-in mux in the way we handle things
		s.list(w, r)
		return
	}
	urlParts := strings.Split(r.URL.Path, "/")
	var err error
	var e *erroneous.Error
	seg := urlParts[len(urlParts)-1]
	guid, err2 := uuid.Parse(seg)
	if err2 != nil {
		err = fmt.Errorf("%s is not a valid guid", seg)
	} else {
		e, err = s.logger.GetError(erroneous.Guid(guid))
	}
	fmt.Println(singleTpl.Execute(w, map[string]interface{}{"Error": e, "Msg": err, "Prefix": s.prefix}))
}

const listTemplate = `
{{define "title"}}{{.Count}} Errors{{end}}
{{define "content"}}
    {{if eq .Count 0}} No Errors Yet!
    {{else}}
        <h2 id="errorcount">ApplicationName - {{.Count}} Errors; last {{.Ago}} ago</h2>
        <table id="ErrorLog" class="alt-rows">
            <thead>
                <tr>
                    <th class="type-col">&nbsp;</th>
                    <th>Id</th>
                    <th>Error</th>
                    <th>Time</th>
                    <th>Url</th>
                    <th>Remote IP</th>
                    <th>Site</th>
                    <th>Server</th>
                </tr>
            </thead>
            <tbody>
            {{range .Errors}}
                <tr>
                    <td></td>
                    <td><a href="{{$.Prefix}}/{{.GUID}}">{{.GUID}}</a></td>
                    <td>{{.Message}}</td>
                    <td>{{.CreationDate.Format "Jan 02, 2006 15:04:05 UTC"}}</td>
                </tr>
            {{end}}
            </tbody>
        </table>
    {{end}}
{{end}}
`
const singleTemplate = `
{{define "title"}}Show Error{{end}}
{{define "content"}}
    {{if .Error}}
        {{.Error.GUID}}
        <br/>{{.Error.Message}}<br/>
        <pre>{{.Error.Detail}}</pre>
    {{end}}
    {{if .Msg}}CRAP! {{.Msg}} {{end}}
{{end}}`

const baseTemplate = `
<html>
<head>
<title>Erroneous - {{block "title" .}}{{end}}</title>
</head>
<body>{{block "content" .}}{{end}}</body>
</html>
`

var listTpl = template.Must(template.New("base").Parse(baseTemplate))
var singleTpl = template.Must(template.New("base").Parse(baseTemplate))

func init() {
	template.Must(listTpl.New("list").Parse(listTemplate))
	template.Must(singleTpl.New("single").Parse(singleTemplate))
}
