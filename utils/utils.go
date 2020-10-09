package utils

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aedavelli/appver-resource/models"
)

const (
	black   = "\033[1;30m"
	red     = "\033[1;31m"
	green   = "\033[1;32m"
	yellow  = "\033[1;33m"
	purple  = "\033[1;34m"
	magenta = "\033[1;35m"
	teal    = "\033[1;36m"
	white   = "\033[1;37m"
	reset   = "\033[0m"
)

var (
	contentType = http.CanonicalHeaderKey("content-type")
)

func printColor(w io.Writer, c string, a ...interface{}) {
	fmt.Fprintf(w, "%s%s%s\n", c, fmt.Sprint(a...), reset)
}

func PrintSuccess(a ...interface{}) {
	printColor(os.Stderr, green, a...)
}

func PrintWarning(a ...interface{}) {
	printColor(os.Stderr, yellow, a...)
}

func Print(a ...interface{}) {
	printColor(os.Stderr, reset, a...)
}

func Fatal(a ...interface{}) {
	printColor(os.Stderr, red, a...)
	os.Exit(1)
}

func GetHttpVersion(s models.Source) *http.Response {
	client := &http.Client{}
	if s.Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	httpReq, err := http.NewRequest("GET", s.Url, nil)
	if err != nil {
		Fatal(err)
	}

	accept := s.Accept
	if accept == "" {
		accept = "application/json"
	}
	httpReq.Header.Add("Accept", accept)

	if s.UserName != "" || s.Password != "" {
		httpReq.SetBasicAuth(s.UserName, s.Password)
	}

	res, err := client.Do(httpReq)
	if err != nil {
		Fatal(err)
	}

	return res
}

func parseXml(s models.Source, r io.Reader) map[string]string {
	version_field := s.VersionField
	scopes := strings.Split(version_field, "::")
	depth := 0
	stack := make([]string, len(scopes))

	decxml := xml.NewDecoder(r)

	m := map[string]string{}
	key := ""
	for {
		tok, _ := decxml.Token()
		if tok == nil {
			break
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			key = tok.Name.Local
			if depth < len(scopes) &&
				key != scopes[depth] &&
				depth != len(scopes)-1 {
				decxml.Skip()
				break
			}

			depth++
			if depth > len(scopes) {
				depth--
				decxml.Skip()
				break
			}
			stack[depth-1] = key

		case xml.EndElement:
			depth--
		case xml.CharData:
			if depth == len(scopes) && strings.Trim(string(tok), " \n") != "" {
				m[key] = string(tok)
			}
		}
	}
	return m
}

func parseJson(s models.Source, r io.Reader) map[string]string {
	version_field := s.VersionField
	scopes := strings.Split(version_field, "::")

	dec := json.NewDecoder(r)

	m := map[string]string{}
	mj := map[string]interface{}{}
	if err := dec.Decode(&mj); err != nil {
		Fatal(err)
	}

	tm := map[string]interface{}{}
	ok := false
	for i := 0; i < len(scopes)-1; i++ {
		if tm, ok = mj[scopes[i]].(map[string]interface{}); !ok {
			Fatal(version_field, "is not present in the version")
		}
		mj = tm
	}

	for k, v := range mj {
		switch vt := v.(type) {
		case string:
			m[k] = vt
		case float64:
			m[k] = fmt.Sprint(vt)
		case bool:
			m[k] = fmt.Sprint(vt)
		default:
			PrintWarning(vt, "will be ignored")
		}
	}
	return m
}

func parseText(r io.Reader) map[string]string {
	m := map[string]string{}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		Fatal(err)
	}
	m["version"] = string(b)
	return m
}

func ParseResponse(s models.Source, res *http.Response) map[string]string {
	defer res.Body.Close()

	var ctypes = res.Header[http.CanonicalHeaderKey("content-type")]
	var ctype = ""
	if len(ctypes) > 0 {
		ctype = strings.Split(ctypes[0], ";")[0]
	}
	if ctype == "" || s.VersionField == "" {
		ctype = "text/plain"
	}

	switch ctype {
	case "application/json":
		return parseJson(s, res.Body)
	case "application/xml":
		return parseXml(s, res.Body)
	case "text/plain":
		return parseText(res.Body)
	default:
		Fatal("Unknown content type", ctype)
	}
	return map[string]string{}
}
