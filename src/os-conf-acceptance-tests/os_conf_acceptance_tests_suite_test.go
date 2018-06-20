package os_conf_acceptance_tests_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	boshBinaryPath string
	deploymentName string
	boshStemcell   string
)

func TestOsConfAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OsConfAcceptanceTests Suite")
}

var _ = BeforeSuite(func() {
	var found bool
	boshBinaryPath = assertEnvExists("BOSH_BINARY_PATH")
	deploymentName = assertEnvExists("BOSH_DEPLOYMENT")
	boshStemcell, found = os.LookupEnv("BOSH_STEMCELL")
	if !found {
		boshStemcell = "ubuntu-trusty"
	}

	assertEnvExists("BOSH_CLIENT")
	assertEnvExists("BOSH_CLIENT_SECRET")
	assertEnvExists("BOSH_CA_CERT")
	assertEnvExists("BOSH_ENVIRONMENT")

	deployOSConfDeployment()
})

var _ = AfterSuite(func() {
	destroyOSConfDeployment()
})

func assertEnvExists(env string) string {
	val, exists := os.LookupEnv(env)
	if !exists {
		Fail(fmt.Sprintf("Expected %s", env))
	}

	return val
}

func deployOSConfDeployment() {
	cmd := exec.Command(
		boshBinaryPath,
		"-n",
		"-d",
		deploymentName,
		"deploy",
		"assets/manifest.yml",
		"-v",
		fmt.Sprintf("stemcell_os=%s", boshStemcell),
		"-v",
		fmt.Sprintf("deployment_name=%s", deploymentName),
	)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
}

func destroyOSConfDeployment() {
	cmd := exec.Command(
		boshBinaryPath,
		"-n",
		"-d",
		deploymentName,
		"delete-deployment",
	)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
}

func boshSSH(jobName, command string) *gexec.Session {
	cmd := exec.Command(
		boshBinaryPath,
		"-d",
		deploymentName,
		"ssh",
		jobName,
		"-c",
		command,
	)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
