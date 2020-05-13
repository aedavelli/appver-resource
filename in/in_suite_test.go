package main_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/aedavelli/appver-resource/test_helper"
)

const port = 8383

var inPath string

var _ = BeforeSuite(func() {
	var err error

	if _, err = os.Stat("/opt/resource/in"); err == nil {
		inPath = "/opt/resource/in"
	} else {
		inPath, err = gexec.Build("github.com/aedavelli/appver-resource/in")
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestIn(t *testing.T) {
	go test_helper.RunTmpServer(port)
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}
