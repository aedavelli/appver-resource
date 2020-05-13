package test_helper

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aedavelli/appver-resource/utils"
)

type KV struct {
	Key   string `json:"key" xml:"key"`
	Value string `json:"value" xml:"value"`
}

type TVersion struct {
	Version     string `json:"version" xml:"version"`
	Repo        string `json:"repo" xml:"repo"`
	Commit      string `json:"commit" xml:"commit"`
	Name        string `json:"name" xml:"name"`
	KVs         []KV   `json:"kvs" xml:"kvs"`
	ChangeCount int    `json:"change_count" xml:"change_count"`
}

type TVersionInfo struct {
	Git TVersion `json:"git" xml:"git"`
	Hg  TVersion `json:"hg" xml:"hg"`
	Svn TVersion `json:"svn" xml:"svn"`
}

type Info struct {
	VersionInfo TVersionInfo `json:"version_info" xml:"version_info"`
}

var i Info

func init() {
	vg := TVersion{
		Version:     "1.0.0_git",
		Repo:        "git@nohost.com:norep",
		Commit:      "d752151f566d70e85f10c5c0b9b1911ab1ed22d2",
		Name:        "Change 1",
		ChangeCount: 222,
	}
	vh := TVersion{
		Version:     "1.0.0_hg",
		Repo:        "hg@nohost.com:norep",
		Commit:      "d752151f566d70e85f10c5c0b9b1911ab1ed22d2",
		Name:        "Change 1",
		ChangeCount: 222,
	}
	vs := TVersion{
		Version:     "1.0.0_svn",
		Repo:        "svn@nohost.com:norep",
		Commit:      "d752151f566d70e85f10c5c0b9b1911ab1ed22d2",
		Name:        "Change 1",
		ChangeCount: 222,
	}
	vg.KVs = []KV{{"kk1", "vv1"}, {"kk2", "vv2"}}
	vh.KVs = []KV{{"kk1", "vv1"}, {"kk2", "vv2"}}
	vs.KVs = []KV{{"kk1", "vv1"}, {"kk2", "vv2"}}
	vi := TVersionInfo{vg, vh, vs}
	i = Info{vi}
}

func versionHandler(w http.ResponseWriter, req *http.Request) {
	var b []byte
	var err error
	var accept string

	if len(req.Header["Accept"]) > 0 {
		accept = req.Header["Accept"][0]
	}

	switch strings.ToLower(accept) {
	case "application/xml":
		w.Header().Add("Content-Type", "application/xml")
		w.Write([]byte("\n"))
		b, err = xml.MarshalIndent(i, "  ", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
	case "application/json":
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("\n"))
		b, err = json.MarshalIndent(i, "  ", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
	default:
		w.Header().Add("Content-Type", "text/plain")
		b = []byte(i.VersionInfo.Git.Version)
	}
	w.Write(b)
}

func RunTmpServer(port int16) {
	utils.PrintWarning("Starting test server")
	http.HandleFunc("/version", versionHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
