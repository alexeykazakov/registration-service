package link

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/codeready-toolchain/registration-service/pkg/signup"
)

// Link represents a link payload
type Link struct {
	OutputURL string
}

type LinkType interface {
	Matches(inputPath string) bool
	OutputURL(inputURL *url.URL, signup signup.Signup) (string, error)
}

type linkTypeImpl struct {
	inputPath      string
	outputTemplate *template.Template
	getParams      func(inputURL *url.URL, signup signup.Signup) params
}

type params struct {
	ConsoleMember   string
	UserNamespace   string
	Path            string
	RHODSMember     string
	DevSpacesMember string
}

var openShiftAdd = &linkTypeImpl{
	inputPath:      "/openshift/add/",
	outputTemplate: newTemplate("add", "{{.ConsoleMember}}/add/ns/{{.UserNamespace}}"),
	getParams: func(_ *url.URL, signup signup.Signup) params {
		return params{
			ConsoleMember: getConsoleMember(signup),
			UserNamespace: getUserNamespace(signup),
		}
	},
}

var rhods = &linkTypeImpl{
	inputPath:      "/datascience/",
	outputTemplate: newTemplate("rhods", "{{.RHODSMember}}/notebookController/spawner"),
	getParams: func(_ *url.URL, signup signup.Signup) params {
		return params{
			RHODSMember: getRHODSMember(signup),
		}
	},
}

var devSpaces = &linkTypeImpl{
	inputPath:      "/devspaces/",
	outputTemplate: newTemplate("che", "{{.DevSpacesMember}}/"),
	getParams: func(_ *url.URL, signup signup.Signup) params {
		return params{
			DevSpacesMember: getDevSpacesMember(signup),
		}
	},
}

var webConsoleBookmark = &linkTypeImpl{
	inputPath:      "/k8s/",
	outputTemplate: newTemplate("console", "{{.ConsoleMember}}/{{.Path}}"),
	getParams: func(inputURL *url.URL, signup signup.Signup) params {
		return params{
			ConsoleMember: getConsoleMember(signup),
			Path:          inputURL.Path,
		}
	},
}

var DefaultLink = &linkTypeImpl{
	outputTemplate: newTemplate("default", "{{.ConsoleMember}}/topology/ns/{{.UserNamespace}}"),
	getParams: func(_ *url.URL, signup signup.Signup) params {
		return params{
			ConsoleMember: getConsoleMember(signup),
			UserNamespace: getUserNamespace(signup),
		}
	},
}

var LinkTypes = []LinkType{openShiftAdd, rhods, devSpaces, webConsoleBookmark}

func newTemplate(name, outputTemplate string) *template.Template {
	tmpl, err := template.New(name).Parse(outputTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func (l *linkTypeImpl) OutputURL(inputURL *url.URL, signup signup.Signup) (string, error) {
	var result bytes.Buffer

	if err := l.outputTemplate.Execute(&result, l.getParams(inputURL, signup)); err != nil {
		return "", err
	}

	return result.String(), nil
}

func (l *linkTypeImpl) Matches(inputPath string) bool {
	if !strings.HasSuffix(inputPath, "/") {
		inputPath = inputPath + "/"
	}

	return strings.HasPrefix(inputPath, l.inputPath)
}

func getConsoleMember(signup signup.Signup) string {
	cUrl := signup.ConsoleURL
	if strings.HasSuffix(cUrl, "/") {
		cUrl = cUrl[:len(cUrl)-1]
	}
	return cUrl
}

func getRHODSMember(signup signup.Signup) string {
	return getAppsURL(signup, "rhods-dashboard-redhat-ods-applications")
}

func getDevSpacesMember(signup signup.Signup) string {
	return getAppsURL(signup, "devspaces")
}

func getUserNamespace(signup signup.Signup) string {
	return fmt.Sprintf("%s-dev", signup.CompliantUsername)
}

// getAppsURL returns a URL for the specific app
// for example for the "devspaces" app and api server "https://api.host.openshiftapps.com:6443"
// it will return "https://devspaces.apps.host.openshiftapps.com"
func getAppsURL(signup signup.Signup, appRouteName string) string {
	return fmt.Sprintf("https://%s.apps.%s", appRouteName, signup.ClusterName)
}
