package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

const (
	defaultUsernameEnvVar = "GIT_USERNAME"
	defaultPasswordEnvVar = "GIT_PASSWORD"
)

// GitOptions specifies the options to use when cloning a git repo and checking out a ref.
type GitOptions struct {
	// RepoURL is the URL of the git repo to clone.
	RepoURL string

	// Ref is the git ref to checkout after cloning.
	Ref string

	// UsernameEnvVar is the name of the environment variable to lookup for the username to use when authing to the git
	// repo. If unset, defaults to defaultUsernameEnvVar.
	UsernameEnvVar string

	// PasswordEnvVar is the name of the environment variable to lookup for the password to use when authing to the git
	// repo. If unset, defaults to defaultPasswordEnvVar.
	PasswordEnvVar string
}

// CloneAndCheckout uses go-git to clone the configured repository at the desired ref into the given cloneDir.
func CloneAndCheckout(opts GitOptions, cloneDir string) error {
	cloneOpts := git.CloneOptions{
		URL:  opts.RepoURL,
		Auth: getGitAuth(opts),
	}

	// If the ref is not a commit hash, assume branch or tag. Note that to checkout a hash using go-git, we need to
	// first clone the repo at HEAD, and then fetch the commits in the repo so that the reference exists. Given that, we
	// handle the commit hash refs after the repo is cloned.
	refIsHash := plumbing.IsHash(opts.Ref)
	if !refIsHash {
		ref, err := getGitReferenceName(opts, cloneDir)
		if err != nil {
			return err
		}
		cloneOpts.ReferenceName = ref
	}

	repo, err := git.PlainClone(cloneDir, false, &cloneOpts)
	if err != nil {
		return err
	}

	// Handle commit hash ref by fetching the references and then checking out the hash directly.
	if refIsHash {
		err := repo.Fetch(&git.FetchOptions{
			RemoteName: "origin",
			Auth:       cloneOpts.Auth,
			RefSpecs:   []config.RefSpec{config.RefSpec(opts.Ref + ":" + opts.Ref)},
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}

		// ... retrieving the commit being pointed by HEAD
		if _, err := repo.Head(); err != nil {
			return err
		}

		tree, err := repo.Worktree()
		if err != nil {
			return err
		}

		// ... checking out to desired commit
		checkoutOpts := &git.CheckoutOptions{Hash: plumbing.NewHash(opts.Ref)}
		if err = tree.Checkout(checkoutOpts); err != nil {
			return err
		}
	}
	return nil
}

// getGitReferenceName converts the string based ref into a valid reference that go-git understands. To do this, this
// routine queries the git repo for the list of references and matches the provided string ref against it. This will
// return an error if the provided ref string does not exist in the git repo.
func getGitReferenceName(opts GitOptions, cloneDir string) (plumbing.ReferenceName, error) {
	remote := git.NewRemote(
		filesystem.NewStorage(
			osfs.New(cloneDir),
			cache.NewObjectLRUDefault(),
		),
		&config.RemoteConfig{
			URLs: []string{opts.RepoURL},
		},
	)

	auth := getGitAuth(opts)
	allRefs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return plumbing.HEAD, err
	}

	if ref := plumbing.NewBranchReferenceName(opts.Ref); gitRefExists(ref, allRefs) {
		return ref, nil
	}

	if ref := plumbing.NewTagReferenceName(opts.Ref); gitRefExists(ref, allRefs) {
		return ref, nil
	}

	return plumbing.HEAD, fmt.Errorf("invalid ref: %s", opts.Ref)
}

// gitRefExists returns true if the given query ref is in the list of all refs.
func gitRefExists(query plumbing.ReferenceName, allRefs []*plumbing.Reference) bool {
	for _, ref := range allRefs {
		if query.String() == ref.Name().String() {
			return true
		}
	}
	return false
}

// getGitAuth constructs authentication parameters for go-git, looking up the values from environment variables.
func getGitAuth(opts GitOptions) transport.AuthMethod {
	usernameEnvVar := defaultUsernameEnvVar
	if opts.UsernameEnvVar != "" {
		usernameEnvVar = opts.UsernameEnvVar
	}
	username := os.Getenv(usernameEnvVar)

	passwordEnvVar := defaultPasswordEnvVar
	if opts.PasswordEnvVar != "" {
		passwordEnvVar = opts.PasswordEnvVar
	}
	password := os.Getenv(passwordEnvVar)

	if username != "" || password != "" {
		return &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}
	return nil
}
