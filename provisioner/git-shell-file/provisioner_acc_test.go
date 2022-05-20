// Run with: PACKER_ACC=1 go test -count 1 -v ./provisioner/git-shell-file/provisioner_acc_test.go  -timeout=120m
package git_shell_file

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/environment"

	"github.com/yorinasub17/packer-plugin-git-shell/provisioner/common"
)

//go:embed test-fixtures/public/single_file.pkr.hcl
var testGitFilePublicSingleFile string

//go:embed test-fixtures/public/multiple_file.pkr.hcl
var testGitFilePublicMultipleFile string

//go:embed test-fixtures/public/folder.pkr.hcl
var testGitFilePublicFolder string

//go:embed test-fixtures/private/template.pkr.hcl
var testGitFilePrivate string

func TestAccGitFilePublicSingleFile(t *testing.T) {
	common.RunAccTest(
		t,
		"git_file_provisioner_public_single_file_test",
		testGitFilePublicSingleFile,
		"git-shell-file",
		[]string{
			fmt.Sprintf("docker.example: Cloning repo %s at ref main", common.PublicAutomatedTestingRepoURL),
			"docker.example: Hello",
		},
	)
}

func TestAccGitFilePublicMultipleFile(t *testing.T) {
	common.RunAccTest(
		t,
		"git_file_provisioner_public_multiple_file_test",
		testGitFilePublicMultipleFile,
		"git-shell-file",
		[]string{
			fmt.Sprintf("docker.example: Cloning repo %s at ref test", common.PublicAutomatedTestingRepoURL),
			// Check both file outputs
			"docker.example: Hello.*(\n.*)+docker.example: Hello World from test",
		},
	)
}

func TestAccGitFilePublicFolder(t *testing.T) {
	common.RunAccTest(
		t,
		"git_file_provisioner_public_folder_test",
		testGitFilePublicFolder,
		"git-shell-file",
		[]string{
			fmt.Sprintf("docker.example: Cloning repo %s at ref test", common.PublicAutomatedTestingRepoURL),
			// Check both file outputs
			"docker.example: Hello.*(\n.*)+docker.example: Hello World from test",
		},
	)
}

func TestAccGitFilePrivate(t *testing.T) {
	environment.RequireEnvVar(t, common.TestGitUsernameEnvVar)
	environment.RequireEnvVar(t, common.TestGitTokenEnvVar)

	common.RunAccTest(
		t,
		"git_file_provisioner_private_test",
		testGitFilePrivate,
		"git-shell-file",
		[]string{
			fmt.Sprintf("docker.example: Cloning repo %s at ref main", common.PrivateAutomatedTestingRepoURL),
			`docker.example: Hello \(private\)!`,
		},
	)
}
