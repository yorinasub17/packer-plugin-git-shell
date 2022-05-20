// Run with: PACKER_ACC=1 go test -count 1 -v ./provisioner/git-shell/provisioner_acc_test.go  -timeout=120m
package git_shell

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/environment"

	"github.com/yorinasub17/packer-plugin-git-shell/provisioner/common"
)

//go:embed test-fixtures/public/single_script.pkr.hcl
var testGitShellPublicSingleScript string

//go:embed test-fixtures/public/multiple_script.pkr.hcl
var testGitShellPublicMultipleScript string

//go:embed test-fixtures/private/template.pkr.hcl
var testGitShellPrivate string

func TestAccGitShellPublicSingleScript(t *testing.T) {
	common.RunAccTest(
		t,
		"git_shell_provisioner_public_single_script_test",
		testGitShellPublicSingleScript,
		"git-shell",
		[]string{
			fmt.Sprintf("docker.example: Cloning repo %s at ref main", common.PublicAutomatedTestingRepoURL),
			"docker.example: from public: hello world",
		},
	)
}

func TestAccGitShellPublicMultipleScript(t *testing.T) {
	common.RunAccTest(
		t,
		"git_shell_provisioner_public_multiple_script_test",
		testGitShellPublicMultipleScript,
		"git-shell",
		[]string{
			fmt.Sprintf("docker.example: Cloning repo %s at ref test", common.PublicAutomatedTestingRepoURL),
			// Check both script outputs at the same time to verify ordering
			"docker.example: from test branch.*(\n.*)+docker.example: from public: hello world",
		},
	)
}

func TestAccGitShellPrivate(t *testing.T) {
	environment.RequireEnvVar(t, common.TestGitUsernameEnvVar)
	environment.RequireEnvVar(t, common.TestGitTokenEnvVar)

	common.RunAccTest(
		t,
		"git_shell_provisioner_private_test",
		testGitShellPrivate,
		"git-shell",
		[]string{
			fmt.Sprintf("docker.example: Cloning repo %s at ref main", common.PrivateAutomatedTestingRepoURL),
			"docker.example: from private repo",
		},
	)
}
