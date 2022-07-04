package kernels

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// NB(kkourt): I'm using go-git for the tests here, just because I was curious.
// Executing git commands seems like the safest choice, however, which is why
// go-git is not used in the main code (at least for now).

func TestUrlExamples(t *testing.T) {
	for _, ex := range UrlExamples {
		kernelURL, err := ParseURL(ex.URL)
		assert.Nil(t, err)
		assert.Equal(t, kernelURL, ex.expectedKernelURL)
	}
}

func gitAddTestFile(r *git.Repository, gitDir, file string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	filename := filepath.Join(gitDir, file)
	if err := ioutil.WriteFile(filename, []byte(file), 0644); err != nil {
		return err
	}
	if _, err = w.Add(file); err != nil {
		return err
	}

	if _, err = w.Commit(fmt.Sprintf("add %s", file), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@test.test",
			When:  time.Now(),
		},
	}); err != nil {
		return err
	}

	return nil
}

func gitCreateAndCheckoutBranch(repo *git.Repository, name string) error {
	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	refName := plumbing.NewBranchReferenceName(name)
	ref := plumbing.NewHashReference(refName, headRef.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		return err
	}

	// NB: code below will save to .config, which does not seem to be needed
	// if err = repo.CreateBranch(&config.Branch{
	// 	Name:  "test",
	// 	Merge: refName,
	// }); err != nil {
	// 	return err
	// }

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := w.Checkout(&git.CheckoutOptions{
		Branch: refName,
	}); err != nil {
		return err
	}

	return nil
}

// makeTestGitRepo creates a simple test repository with two branches:
// master:
//   - file1
// branch:
//   - file1
//   - file2
func makeTestGitRepo() (string, error) {
	dir, err := os.MkdirTemp("", "git-test-repo-")
	if err != nil {
		return "", err
	}

	repo, err := git.PlainInit(dir, false)
	if err != nil {
		return "", err
	}

	// add a file to the head (normally master branch)
	err = gitAddTestFile(repo, dir, "file1")
	if err != nil {
		return "", err
	}

	// create and checkout a branch
	err = gitCreateAndCheckoutBranch(repo, "branch")
	if err != nil {
		return "", err
	}

	// add a file to "branch"
	err = gitAddTestFile(repo, dir, "file2")
	if err != nil {
		return "", err
	}

	return dir, nil
}

func TestGitFetch(t *testing.T) {
	log := logrus.New()
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}

	gitRepoDir, err := makeTestGitRepo()
	assert.Nil(t, err)
	t.Logf("git test repository: %s\n", gitRepoDir)
	dir, err := os.MkdirTemp("", "git-test-dir-")
	assert.Nil(t, err)
	t.Logf("git test repository: %s\n", dir)

	// src1: only file1 exists
	gu1 := newGitURL(gitRepoDir, "")
	err = gu1.Fetch(context.Background(), log, dir, "src1")
	assert.Nil(t, err)
	isReg, err := regularFileExists(filepath.Join(dir, "src1", "file1"))
	assert.Nil(t, err)
	assert.True(t, isReg)
	isReg, err = regularFileExists(filepath.Join(dir, "src1", "file2"))
	assert.Nil(t, err)
	assert.False(t, isReg)

	// src2: only file1 and file2 exist
	gu2 := newGitURL(gitRepoDir, "branch")
	err = gu2.Fetch(context.Background(), log, dir, "src2")
	assert.Nil(t, err)
	isReg, err = regularFileExists(filepath.Join(dir, "src2", "file1"))
	assert.Nil(t, err)
	assert.True(t, isReg)
	isReg, err = regularFileExists(filepath.Join(dir, "src2", "file2"))
	assert.Nil(t, err)
	assert.True(t, isReg)

	// src1: only file1 exists
	// TODO: modify the repository by removing file1 and test again
	gu1 = newGitURL(gitRepoDir, "")
	err = gu1.Fetch(context.Background(), log, dir, "src1")
	assert.Nil(t, err)
	isReg, err = regularFileExists(filepath.Join(dir, "src1", "file1"))
	assert.Nil(t, err)
	assert.True(t, isReg)
	isReg, err = regularFileExists(filepath.Join(dir, "src1", "file2"))
	assert.Nil(t, err)
	assert.False(t, isReg)
}
