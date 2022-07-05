package kernels

import (
	"context"
	"fmt"

	"github.com/cilium/little-vm-helper/pkg/logcmd"
	"github.com/sirupsen/logrus"
)

type gitAddWorkdirArg struct {
	workDir      string
	bareDir      string
	remoteName   string
	remoteRepo   string
	remoteBranch string
	localBranch  string
}

func gitAddWorkdir(ctx context.Context, log *logrus.Logger, arg *gitAddWorkdirArg) error {
	remoteAddArgs := []string{
		"--git-dir", arg.bareDir,
		"remote", "add",
		"-f", "-t", arg.remoteBranch, arg.remoteName, arg.remoteRepo,
	}
	if err := logcmd.RunAndLogCmdContext(ctx, log, GitBinary, remoteAddArgs...); err != nil {
		return err
	}

	worktreeAddArgs := []string{
		"--git-dir", arg.bareDir,
		"worktree", "add",
		"-b", arg.localBranch,
		"--track",
		arg.workDir,
		fmt.Sprintf("%s/%s", arg.remoteName, arg.remoteBranch),
	}

	return logcmd.RunAndLogCmdContext(ctx, log, GitBinary, worktreeAddArgs...)
}

type gitRemoveWorkdirArg struct {
	workDir     string
	bareDir     string
	remoteName  string
	localBranch string
}

func gitRemoveWorkdir(ctx context.Context, log *logrus.Logger, arg *gitRemoveWorkdirArg) {

	worktreeRemoveArgs := []string{
		"--git-dir", arg.bareDir,
		"worktree", "remove",
		arg.workDir,
	}
	if err := logcmd.RunAndLogCmdContext(ctx, log, GitBinary, worktreeRemoveArgs...); err != nil {
		log.WithError(err).Warn("did not remove worktree")
	}

	remoteRemoveArgs := []string{
		"--git-dir", arg.bareDir,
		"remote", "remove",
		arg.remoteName,
	}
	if err := logcmd.RunAndLogCmdContext(ctx, log, GitBinary, remoteRemoveArgs...); err != nil {
		log.WithError(err).Warn("did not remove remote")
	}

	branchRemoveArgs := []string{
		"--git-dir", arg.bareDir,
		"branch", "--delete", "--force",
		arg.localBranch,
	}
	if err := logcmd.RunAndLogCmdContext(ctx, log, GitBinary, branchRemoveArgs...); err != nil {
		log.WithError(err).Warn("did not remove local branch")
	}
}
