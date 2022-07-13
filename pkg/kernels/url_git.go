package kernels

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cilium/little-vm-helper/pkg/logcmd"
	"github.com/sirupsen/logrus"
)

// NB: we use the same git directory for all remotes so that we can speed up
// downloads.  The name of the kernel as provided by the user is used as the
// remote name. This means that we might end up with the same remote on
// different names, but that's fine.
var MainGitDir = "git"

type GitURL struct {
	Repo string

	// branch in remote (by default, master)
	Branch string
}

func NewGitURL(kurl *url.URL) (KernelURL, error) {
	// NB: far from perfect, but works for the simple cases
	repo := fmt.Sprintf("%s://%s%s", kurl.Scheme, kurl.Host, kurl.Path)
	// NB: we (ab)use the fragment part of the URL to store the branch
	branch := kurl.Fragment
	return newGitURL(repo, branch), nil
}

func newGitURL(repo string, branch string) *GitURL {
	if branch == "" {
		branch = "master"
	}

	return &GitURL{
		Repo:   repo,
		Branch: branch,
	}
}

func (gu *GitURL) syncWorktree(
	ctx context.Context,
	log *logrus.Logger,
	idDir string,
) error {
	oldPath, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(idDir)
	if err != nil {
		return err
	}
	defer os.Chdir(oldPath)

	cmd := exec.CommandContext(ctx, "git", "pull")
	return logcmd.RunAndLogCommand(cmd, log)
}

func makeGitDir(ctx context.Context, log *logrus.Logger, gitDir string) error {
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "git", "init", "--bare", gitDir)
	if err := logcmd.RunAndLogCommand(cmd, log); err != nil {
		os.RemoveAll(gitDir)
		return err
	}

	return nil
}

// fetch will fetches the code pointed by gu, into dir/id
// It uses a
func (gu *GitURL) fetch(
	ctx context.Context,
	log *logrus.Logger,
	dir string,
	id string,
) error {

	if err := CheckEnvironment(); err != nil {
		return err
	}

	if id == MainGitDir {
		return fmt.Errorf("id `%s` is not allowed. Please use another.", id)
	}

	// directories are
	// <dir>/<MainGitDir> ->  git repo
	// <dir>/<id> -> one worktree per id
	gitDir := filepath.Join(dir, MainGitDir)
	idDir := filepath.Join(dir, id)

	if idExists, err := directoryExists(idDir); err != nil {
		return err
	} else if idExists {
		return gu.syncWorktree(ctx, log, idDir)
	}

	if gitExists, err := directoryExists(gitDir); err != nil {
		return err
	} else if !gitExists {
		if err := makeGitDir(ctx, log, gitDir); err != nil {
			return err
		}
	}

	return gitAddWorkdir(ctx, log, &gitAddWorkdirArg{
		workDir:      idDir,
		bareDir:      gitDir,
		remoteName:   id,
		remoteRepo:   gu.Repo,
		remoteBranch: gu.Branch,
		localBranch:  gitLocalBranch(id),
	})
}
