// Run with: PACKER_ACC=1 go test -count 1 -v ./provisioner/git-shell-local/provisioner_acc_test.go  -timeout=120m
package git_shell_local

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/environment"

	"github.com/yorinasub17/packer-plugin-git-shell/provisioner/common"
)

//go:embed test-fixtures/public/single_script.pkr.hcl
var testGitShellLocalPublicSingleScript string

//go:embed test-fixtures/public/multiple_script.pkr.hcl
var testGitShellLocalPublicMultipleScript string

//go:embed test-fixtures/private/template.pkr.hcl
var testGitShellLocalPrivate string

func TestAccGitShellLocalPublicSingleScript(t *testing.T) {
	common.RunAccTest(
		t,
		"git_shell_local_provisioner_public_single_script_test",
		testGitShellLocalPublicSingleScript,
		"git-shell-local",
		[]string{
			fmt.Sprintf("null.basic-example: Cloning repo %s at ref main", common.PublicAutomatedTestingRepoURL),
			"null.basic-example: from public: hello world",
		},
	)
}

func TestAccGitShellLocalPublicMultipleScript(t *testing.T) {
	common.RunAccTest(
		t,
		"git_shell_local_provisioner_public_multiple_script_test",
		testGitShellLocalPublicMultipleScript,
		"git-shell-local",
		[]string{
			fmt.Sprintf("null.basic-example: Cloning repo %s at ref test", common.PublicAutomatedTestingRepoURL),
			// Check both script outputs at the same time to verify ordering
			"null.basic-example: from test branch.*\n.*null.basic-example: from public: hello world",
		},
	)
}

func TestAccGitShellLocalPrivate(t *testing.T) {
	environment.RequireEnvVar(t, common.TestGitUsernameEnvVar)
	environment.RequireEnvVar(t, common.TestGitTokenEnvVar)

	common.RunAccTest(
		t,
		"git_shell_local_provisioner_private_test",
		testGitShellLocalPrivate,
		"git-shell-local",
		[]string{
			fmt.Sprintf("null.basic-example: Cloning repo %s at ref main", common.PrivateAutomatedTestingRepoURL),
			"null.basic-example: from private repo",
		},
	)
}
