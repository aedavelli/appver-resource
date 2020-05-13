package main_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/aedavelli/appver-resource/models"
)

const (
	version_git = "1.0.0_git"
	version_hg  = "1.0.0_hg"
	version_svn = "1.0.0_svn"
	version_old = "0.1.0"
)

var _ = Describe("Check", func() {
	var (
		checkCmd *exec.Cmd
		request  models.CheckRequest
		response models.CheckResponse
	)

	BeforeEach(func() {
		request = models.CheckRequest{
			Source: models.Source{
				Url: fmt.Sprintf("http://127.0.0.1:%d/version", port),
			},
		}

		response = models.CheckResponse{}

		checkCmd = exec.Command(checkPath)

	})

	JustBeforeEach(func() {
		stdin, err := checkCmd.StdinPipe()
		Expect(err).NotTo(HaveOccurred())

		session, err := gexec.Start(checkCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		err = json.NewEncoder(stdin).Encode(request)
		Expect(err).NotTo(HaveOccurred())

		<-session.Exited
		Expect(session.ExitCode()).To(Equal(0))

		err = json.Unmarshal(session.Out.Contents(), &response)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("No Version Supplied", func() {
		BeforeEach(func() {
			request.Source.Accept = "application/json"
		})
		Context("when version control is git", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::git::version"
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(1))
				Expect(response[0].Version).To(Equal(version_git))
			})
		})

		Context("when version control is hg", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::hg::version"
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(1))
				Expect(response[0].Version).To(Equal(version_hg))
			})
		})

		Context("when version control is svn", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::svn::version"
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(1))
				Expect(response[0].Version).To(Equal(version_svn))
			})
		})
	})

	Describe("Version Other than latest Supplied", func() {
		BeforeEach(func() {
			request.Version = models.AppVersion{version_old}
			request.Source.Accept = "application/json"
		})
		Context("when version control is git", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::git::version"
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(2))
				Expect(response[0].Version).To(Equal(version_old))
				Expect(response[1].Version).To(Equal(version_git))
			})
		})

		Context("when version control is hg", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::hg::version"
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(2))
				Expect(response[0].Version).To(Equal(version_old))
				Expect(response[1].Version).To(Equal(version_hg))
			})
		})

		Context("when version control is svn", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::svn::version"
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(2))
				Expect(response[0].Version).To(Equal(version_old))
				Expect(response[1].Version).To(Equal(version_svn))
			})
		})

	})

	Describe("Same Version Supplied", func() {
		BeforeEach(func() {
			request.Source.Accept = "application/json"
		})
		Context("when version control is git", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::git::version"
				request.Version = models.AppVersion{version_git}
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(1))
				Expect(response[0].Version).To(Equal(version_git))
			})
		})

		Context("when version control is hg", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::hg::version"
				request.Version = models.AppVersion{version_hg}
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(1))
				Expect(response[0].Version).To(Equal(version_hg))
			})
		})

		Context("when version control is svn", func() {
			BeforeEach(func() {
				request.Source.VersionField = "version_info::svn::version"
				request.Version = models.AppVersion{version_svn}
			})

			It("Check Version fron response Url and size", func() {
				Expect(response).To(HaveLen(1))
				Expect(response[0].Version).To(Equal(version_svn))
			})
		})
	})
})
