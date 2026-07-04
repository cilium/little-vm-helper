// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"archive/tar"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTarPathCreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()

	tarPath := filepath.Join(dir, "kernel.tar")
	file, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("create tar: %v", err)
	}

	tw := tar.NewWriter(file)
	content := []byte("kernel")
	header := &tar.Header{
		Name: "data/kernels/6.6-main/boot/vmlinuz-6.6.0",
		Mode: 0o644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(header); err != nil {
		t.Fatalf("write header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("write file content: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close tar file: %v", err)
	}

	targetDir := filepath.Join(dir, "out")
	if err := ExtractTarPath(tarPath, "data/kernels/6.6-main", targetDir); err != nil {
		t.Fatalf("extract tar path: %v", err)
	}

	extracted := filepath.Join(targetDir, "boot", "vmlinuz-6.6.0")
	got, err := os.ReadFile(extracted)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("unexpected extracted content: got %q want %q", got, content)
	}
}
