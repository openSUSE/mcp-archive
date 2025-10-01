// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package archive

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestCpioList(t *testing.T) {
	files, err := cpioList("../../testdata/test.cpio", 0)
	if err != nil {
		t.Fatalf("cpioList failed: %v", err)
	}

	expected := []string{
		"foo 0",
		"foo/baar.txt 27",
		"foo/bazz 5",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	// For this specific case, the order is deterministic, but this is a good practice.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			// This is a simplification. A real test might parse the lines and compare fields.
			// For example, the mode can vary based on umask.
			// For now, we'll do a simple substring match on name and size.
			if strings.Contains(file, exp) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file '%s' not found in archive", exp)
		}
	}
}

func TestTarGzList(t *testing.T) {
	files, err := tarGzList("../../testdata/test.tar.gz", 0)
	if err != nil {
		t.Fatalf("tarGzList failed: %v", err)
	}

	expected := []string{
		"foo/ 0",
		"foo/baar.txt 27",
		"foo/bazz 5",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if strings.Contains(file, exp) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file '%s' not found in archive", exp)
		}
	}
}

func TestTarBz2List(t *testing.T) {
	files, err := tarBz2List("../../testdata/test.tar.bz2", 0)
	if err != nil {
		t.Fatalf("tarBz2List failed: %v", err)
	}

	expected := []string{
		"foo/ 0",
		"foo/baar.txt 27",
		"foo/bazz 5",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if strings.Contains(file, exp) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file '%s' not found in archive", exp)
		}
	}
}

func TestTarXzList(t *testing.T) {
	files, err := tarXzList("../../testdata/test.tar.xz", 0)
	if err != nil {
		t.Fatalf("tarXzList failed: %v", err)
	}

	expected := []string{
		"foo/ 0",
		"foo/baar.txt 27",
		"foo/bazz 5",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if strings.Contains(file, exp) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file '%s' not found in archive", exp)
		}
	}
}

func TestZipList(t *testing.T) {
	files, err := zipList("../../testdata/test.zip", 0)
	if err != nil {
		t.Fatalf("zipList failed: %v", err)
	}

	expected := []string{
		"foo/ 0",
		"foo/baar.txt 27",
		"foo/bazz 5",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if strings.Contains(file, exp) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file '%s' not found in archive", exp)
		}
	}
}

func TestExtractArchiveFiles(t *testing.T) {
	// Extract a file from the archive.
	_, result, err := ExtractArchiveFiles(context.Background(), &mcp.CallToolRequest{}, ExtractArchiveFilesArgs{
		Path:  "../../testdata/test.zip",
		Files: []string{"foo/baar.txt"},
	})
	if err != nil {
		t.Fatalf("ExtractArchiveFiles failed: %v", err)
	}
	extractedFiles, ok := result.([]File)
	if !ok {
		t.Fatalf("unexpected result type: %T", result)
	}
	if len(extractedFiles) != 1 {
		t.Fatalf("expected 1 file, got %d", len(extractedFiles))
	}
	file := extractedFiles[0]
	if file.Name != "foo/baar.txt" {
		t.Errorf("unexpected file name: %s", file.Name)
	}
	if file.Content != "das Pferd isst Gurkensalat\n" {
		t.Errorf("unexpected content in extracted file: %s", file.Content)
	}
	if file.Size != 27 {
		t.Errorf("unexpected file size: %d", file.Size)
	}
}

func TestExtractArchiveFiles_SizeLimit(t *testing.T) {
	// Set a small size limit for this test.
	originalSizeLimit := MaxExtractFileSize
	MaxExtractFileSize = 20
	defer func() { MaxExtractFileSize = originalSizeLimit }()
	// Attempt to extract the file which is larger than the limit.
	_, _, err := ExtractArchiveFiles(context.Background(), &mcp.CallToolRequest{}, ExtractArchiveFilesArgs{
		Path:  "../../testdata/test.zip",
		Files: []string{"foo/baar.txt"},
	})
	if err == nil {
		t.Fatal("expected error for large file, but got nil")
	}
	if !strings.Contains(err.Error(), "is too large") {
		t.Fatalf("expected size limit error, got: %v", err)
	}
}
