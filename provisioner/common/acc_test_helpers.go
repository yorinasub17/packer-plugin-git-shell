package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/stretchr/testify/assert"
)

// RunAccTest is a helper routine to simplify setting up an acc test for the various provisioners in this repo.
func RunAccTest(
	t *testing.T,
	testName string,
	testTemplate string,
	provisionerType string,
	expectedLogs []string,
) {
	testCase := &acctest.PluginTestCase{
		Name: testName,
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testTemplate,
		Type:     provisionerType,
		Init:     true,
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := ioutil.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)
			for _, expectedLog := range expectedLogs {
				assert.Regexp(t, regexp.MustCompile(expectedLog+".*"), logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
