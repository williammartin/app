package integration_test

import (
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var runaBinPath string

var _ = SynchronizedBeforeSuite(func() []byte {
	preloadAssetImageCmd := exec.Command("docker", "pull", "cfgarden/hello")
	Eventually(execBin(preloadAssetImageCmd), time.Hour).Should(gexec.Exit(0))

	var err error
	runaBinPath, err = gexec.Build("github.com/williammartin/runa")
	Expect(err).NotTo(HaveOccurred())

	return []byte{}
}, func(_ []byte) {})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func execBin(cmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
