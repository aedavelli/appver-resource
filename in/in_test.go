package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/aedavelli/appver-resource/models"
)

const (
	version_git  = "1.0.0_git"
	repo_git     = "git@nohost.com:norep"
	version_hg   = "1.0.0_hg"
	repo_hg      = "hg@nohost.com:norep"
	version_svn  = "1.0.0_svn"
	repo_svn     = "svn@nohost.com:norep"
	commit       = "d752151f566d70e85f10c5c0b9b1911ab1ed22d2"
	change_count = "222"
	name         = "Change 1"
)

var _ = Describe("In", func() {
	var (
		tmpdir      string
		destination string
		inCmd       *exec.Cmd
		request     models.InRequest
		response    models.InResponse
	)

	BeforeEach(func() {
		request = models.InRequest{
			Version: models.AppVersion{
				Version: "1",
			},
			Source: models.Source{
				Url: fmt.Sprintf("http://127.0.0.1:%d/version", port),
			},
		}

		response = models.InResponse{}

		var err error

		tmpdir, err = ioutil.TempDir("", "in-destination")
		Expect(err).NotTo(HaveOccurred())

		destination = path.Join(tmpdir, "in-dir")

		inCmd = exec.Command(inPath, destination)

	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	JustBeforeEach(func() {
		stdin, err := inCmd.StdinPipe()
		Expect(err).NotTo(HaveOccurred())

		session, err := gexec.Start(inCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		err = json.NewEncoder(stdin).Encode(request)
		Expect(err).NotTo(HaveOccurred())

		<-session.Exited
		Expect(session.ExitCode()).To(Equal(0))

		err = json.Unmarshal(session.Out.Contents(), &response)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("JSON Type Payload", func() {
		BeforeEach(func() {
			request.Source.Accept = "application/json"
		})
		Context("when version control is git", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::git::version"
			})

			It("reports the version to be the input version", func() {
				Expect(response.Version.Version).To(Equal("1"))
			})

			It("writes the requested data the destination", func() {
				checkProp(destination, "version", version_git, response.Metadata)
				checkProp(destination, "change_count", change_count, response.Metadata)
				checkProp(destination, "commit", commit, response.Metadata)
				checkProp(destination, "repo", repo_git, response.Metadata)
				checkProp(destination, "name", name, response.Metadata)
			})

		})

		Context("when version control is hg", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::hg::version"
			})

			It("writes the requested data the destination", func() {
				checkProp(destination, "version", version_hg, response.Metadata)
				checkProp(destination, "change_count", change_count, response.Metadata)
				checkProp(destination, "commit", commit, response.Metadata)
				checkProp(destination, "repo", repo_hg, response.Metadata)
				checkProp(destination, "name", name, response.Metadata)
			})

		})

		Context("when version control is svn", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::svn::version"
			})

			It("writes the requested data the destination", func() {
				checkProp(destination, "version", version_svn, response.Metadata)
				checkProp(destination, "change_count", change_count, response.Metadata)
				checkProp(destination, "commit", commit, response.Metadata)
				checkProp(destination, "repo", repo_svn, response.Metadata)
				checkProp(destination, "name", name, response.Metadata)
			})

		})
	})

	Describe("XML Type Payload", func() {
		BeforeEach(func() {
			request.Source.Accept = "application/xml"
		})

		Context("when version control is git", func() {
			BeforeEach(func() {
				request.Source.VersionField = "Info::version_info::git::version"
			})

			It("reports the version to be the input version", func() {
				Expect(response.Version.Version).To(Equal("1"))
			})

			It("writes the requested data the destination", func() {
				checkProp(destination, "version", version_git, response.Metadata)
				checkProp(destination, "change_count", change_count, response.Metadata)
				checkProp(destination, "commit", commit, response.Metadata)
				checkProp(destination, "repo", repo_git, response.Metadata)
				checkProp(destination, "name", name, response.Metadata)
			})

		})

		Context("when version control is hg", func() {
			BeforeEach(func() {
				request.Source.VersionField = "Info::version_info::hg::version"
			})

			It("writes the requested data the destination", func() {
				checkProp(destination, "version", version_hg, response.Metadata)
				checkProp(destination, "change_count", change_count, response.Metadata)
				checkProp(destination, "commit", commit, response.Metadata)
				checkProp(destination, "repo", repo_hg, response.Metadata)
				checkProp(destination, "name", name, response.Metadata)
			})

		})

		Context("when version control is svn", func() {
			BeforeEach(func() {
				request.Source.VersionField = "Info::version_info::svn::version"
			})

			It("writes the requested data the destination", func() {
				checkProp(destination, "version", version_svn, response.Metadata)
				checkProp(destination, "change_count", change_count, response.Metadata)
				checkProp(destination, "commit", commit, response.Metadata)
				checkProp(destination, "repo", repo_svn, response.Metadata)
				checkProp(destination, "name", name, response.Metadata)
			})

		})
	})

	Describe("TEXT Type Payload", func() {
		BeforeEach(func() {
			request.Source.Accept = "text/plain"
		})

		Context("when executed", func() {
			It("reports the version to be the input version", func() {
				Expect(response.Version.Version).To(Equal("1"))
			})

			It("writes the requested data the destination", func() {
				checkProp(destination, "version", version_git, response.Metadata)
				Expect(response.Metadata).To(HaveLen(1))
			})

		})
	})
})

func checkProp(destination, prop, valueToCheck string, meta models.Metadata) {
	output := filepath.Join(destination, prop)
	file, err := ioutil.ReadFile(output)
	Expect(err).NotTo(HaveOccurred())
	val := string(file)
	Expect(val).To(Equal(valueToCheck))
	found := false
	for _, v := range meta {
		if v.Name == prop {
			found = true
			Expect(v.Value).To(Equal(valueToCheck))
		}
	}

	if !found {
		Fail(fmt.Sprintf("%s not found in metadata", prop))
	}
}
