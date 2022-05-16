package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneAndCheckout(t *testing.T) {
	t.Parallel()

	environment.RequireEnvVar(t, TestGitTokenEnvVar)

	testGitPasswordEnvVar := "TEST_GIT_PASSWORD_" + strings.ToUpper(random.UniqueId())
	defer os.Unsetenv(testGitPasswordEnvVar)
	os.Setenv(testGitPasswordEnvVar, os.Getenv(TestGitTokenEnvVar))

	testGitUsernameEnvVar := "TEST_GIT_USERNAME_" + strings.ToUpper(random.UniqueId())
	defer os.Unsetenv(testGitUsernameEnvVar)
	os.Setenv(testGitUsernameEnvVar, "git")

	tests := []struct {
		testName    string
		opts        GitOptions
		script      string
		expectedOut string
		expectErr   bool
	}{
		{
			"public-main",
			GitOptions{
				RepoURL: PublicAutomatedTestingRepoURL,
				Ref:     "main",
			},
			"scripts/echo-sha",
			"17c911e08b4914885e933e235c9c0f78dd8f31b4",
			false,
		},
		{
			"public-main-sha",
			GitOptions{
				RepoURL: PublicAutomatedTestingRepoURL,
				Ref:     "3fc73aab4c9f2d5ade11b665a027b6b830c4862b",
			},
			"scripts/echo-sha",
			"3fc73aab4c9f2d5ade11b665a027b6b830c4862b",
			false,
		},
		{
			"public-branch",
			GitOptions{
				RepoURL: PublicAutomatedTestingRepoURL,
				Ref:     "test",
			},
			"scripts/echo-test",
			"from test branch",
			false,
		},
		{
			"public-branch-sha",
			GitOptions{
				RepoURL: PublicAutomatedTestingRepoURL,
				Ref:     "4900efa4df503e54994064049aa1123a05a66a6f",
			},
			"scripts/echo-test",
			"from test branch",
			false,
		},
		{
			"public-tag",
			GitOptions{
				RepoURL: PublicAutomatedTestingRepoURL,
				Ref:     "v0.0.0",
			},
			"scripts/echo-release",
			"from release tag",
			false,
		},
		{
			"public-no-ref",
			GitOptions{
				RepoURL: PublicAutomatedTestingRepoURL,
				Ref:     "i-dont-exist",
			},
			"scripts/echo-sha",
			"",
			true,
		},
		{
			"private",
			GitOptions{
				RepoURL:        PrivateAutomatedTestingRepoURL,
				Ref:            "main",
				UsernameEnvVar: testGitUsernameEnvVar,
				PasswordEnvVar: testGitPasswordEnvVar,
			},
			"scripts/echo-private",
			"from private repo",
			false,
		},
		{
			"private-no-auth",
			GitOptions{
				RepoURL: PrivateAutomatedTestingRepoURL,
				Ref:     "main",
			},
			"scripts/echo-private",
			"",
			true,
		},
	}

	// Group all parallel tests in a synchronous subtest so that the cleanup routine doesn't run until all subtests are
	// done.
	t.Run("group", func(t *testing.T) {
		for _, test := range tests {
			// Capture range bring into for loop closure so it doesn't change when switching goroutines.
			test := test

			t.Run(test.testName, func(t *testing.T) {
				t.Parallel()

				cloneDir, tmpErr := ioutil.TempDir("", "")
				require.NoError(t, tmpErr)

				err := CloneAndCheckout(test.opts, cloneDir)
				if test.expectErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)

				scriptPath := filepath.Join(cloneDir, test.script)
				cmd := shell.Command{Command: scriptPath, WorkingDir: cloneDir}
				out := shell.RunCommandAndGetStdOut(t, cmd)
				assert.Equal(t, test.expectedOut, out)
			})
		}
	})
}
