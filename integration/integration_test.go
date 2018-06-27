package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration", func() {
	var fixture string

	BeforeEach(func() {
		fixture = "./test_assets/itchy"
	})

	JustBeforeEach(func() {
		appCmd := exec.Command(appBinPath, "build", fixture, "-t", "my/app")
		session := execBin(appCmd)
		Eventually(session).Should(gexec.Exit(0))
	})

	AfterEach(func() {
		removeImageCmd := exec.Command("docker", "image", "rm", "my/app")
		Eventually(execBin(removeImageCmd)).Should(gexec.Exit(0))
	})

	Describe("with a local directory", func() {
		It("creates a docker image with the requested tag", func() {
			dockerCmd := execBin(exec.Command("docker", "image", "list"))
			Eventually(dockerCmd).Should(gbytes.Say("my/app"))
		})

		Describe("the created image", func() {
			var tempDir string

			BeforeEach(func() {
				tempDir = mktmp()
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tempDir)).To(Succeed())
			})

			PIt("is based on the requested rootfs", func() {
			})

			It("has the user's code at the requested location", func() {
				export("my/app", tempDir)
				Expect(filepath.Join(tempDir, "layer0", "tmp", "app", "myfile")).To(BeAnExistingFile())
			})
		})
	})
})

func export(tag, dest string) {
	Eventually(execBin(exec.Command("docker", "save", tag, "-o", filepath.Join(dest, "app.tar")))).Should(gexec.Exit(0))

	tarCmd := exec.Command("tar", "xf", "app.tar")
	tarCmd.Dir = dest
	Eventually(execBin(tarCmd)).Should(gexec.Exit(0))

	s := []struct{ Layers []string }{}
	f, err := os.Open(filepath.Join(dest, "manifest.json"))
	Expect(err).NotTo(HaveOccurred())
	json.NewDecoder(f).Decode(&s)

	for i, layer := range s[0].Layers {
		d := filepath.Join(dest, fmt.Sprintf("layer%d", i))
		Expect(os.MkdirAll(d, 0755)).To(Succeed())
		Eventually(execBin(exec.Command("tar", "xf", filepath.Join(dest, layer), "-C", d))).Should(gexec.Exit(0))
	}
}

func mktmp() string {
	tmp, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())
	return tmp
}
