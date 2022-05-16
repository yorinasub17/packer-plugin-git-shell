// Run with: PACKER_ACC=1 go test -count 1 -v ./provisioner/git-shell-local/provisioner_acc_test.go  -timeout=120m
package git_shell_local

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/stretchr/testify/assert"

	"github.com/yorinasub17/packer-plugin-git-shell/provisioner/common"
)

//go:embed test-fixtures/public/single_script.pkr.hcl
var testGitShellLocalPublicSingleScript string

//go:embed test-fixtures/public/multiple_script.pkr.hcl
var testGitShellLocalPublicMultipleScript string

//go:embed test-fixtures/private/template.pkr.hcl
var testGitShellLocalPrivate string

func TestAccGitShellLocalPublicSingleScript(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "git_shell_local_provisioner_public_single_script_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testGitShellLocalPublicSingleScript,
		Type:     "git-shell-local",
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

			expectedLogs := []string{
				fmt.Sprintf("null.basic-example: Cloning repo %s at ref main", common.PublicAutomatedTestingRepoURL),
				"null.basic-example: from public: hello world",
			}
			for _, expectedLog := range expectedLogs {
				assert.Regexp(t, regexp.MustCompile(expectedLog+".*"), logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccGitShellLocalPublicMultipleScript(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "git_shell_local_provisioner_public_multiple_script_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testGitShellLocalPublicMultipleScript,
		Type:     "git-shell-local",
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

			expectedLogs := []string{
				fmt.Sprintf("null.basic-example: Cloning repo %s at ref test", common.PublicAutomatedTestingRepoURL),
				// Check both script outputs at the same time to verify ordering
				"null.basic-example: from test branch.*\n.*null.basic-example: from public: hello world",
			}
			for _, expectedLog := range expectedLogs {
				assert.Regexp(t, regexp.MustCompile(expectedLog+".*"), logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccGitShellLocalPrivate(t *testing.T) {
	environment.RequireEnvVar(t, common.TestGitUsernameEnvVar)
	environment.RequireEnvVar(t, common.TestGitTokenEnvVar)

	testCase := &acctest.PluginTestCase{
		Name: "git_shell_local_provisioner_private_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testGitShellLocalPrivate,
		Type:     "git-shell-local",
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

			expectedLogs := []string{
				fmt.Sprintf("null.basic-example: Cloning repo %s at ref main", common.PrivateAutomatedTestingRepoURL),
				"null.basic-example: from private repo",
			}
			for _, expectedLog := range expectedLogs {
				assert.Regexp(t, regexp.MustCompile(expectedLog+".*"), logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
