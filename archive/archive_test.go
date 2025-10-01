// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package archive

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newTestArchive(t *testing.T) *Archive {
	a, err := New("../testdata")
	if err != nil {
		t.Fatalf("failed to create archive: %v", err)
	}
	return a
}

type expectedFile struct {
	name string
	size int64
}

func containsFile(files []FileInfo, expected expectedFile) bool {
	for _, file := range files {
		if file.Name == expected.name && file.Size == expected.size {
			return true
		}
	}
	return false
}

func TestCpioList(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.cpioList(filepath.Join(a.Workdir, "test.cpio"), 0)
	if err != nil {
		t.Fatalf("cpioList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo", size: 0},
		{name: "foo/baar.txt", size: 27},
		{name: "foo/bazz", size: 5},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestTarGzList(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.tarGzList(filepath.Join(a.Workdir, "test.tar.gz"), 0)
	if err != nil {
		t.Fatalf("tarGzList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
		{name: "foo/baar.txt", size: 27},
		{name: "foo/bazz", size: 5},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestTarBz2List(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.tarBz2List(filepath.Join(a.Workdir, "test.tar.bz2"), 0)
	if err != nil {
		t.Fatalf("tarBz2List failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
		{name: "foo/baar.txt", size: 27},
		{name: "foo/bazz", size: 5},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestTarXzList(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.tarXzList(filepath.Join(a.Workdir, "test.tar.xz"), 0)
	if err != nil {
		t.Fatalf("tarXzList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
		{name: "foo/baar.txt", size: 27},
		{name: "foo/bazz", size: 5},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestZipList(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.zipList(filepath.Join(a.Workdir, "test.zip"), 0)
	if err != nil {
		t.Fatalf("zipList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
		{name: "foo/baar.txt", size: 27},
		{name: "foo/bazz", size: 5},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestCpioExtract(t *testing.T) {
	a := newTestArchive(t)
	extractedFiles, err := a.cpioExtract(filepath.Join(a.Workdir, "test.cpio"), []string{"foo/baar.txt"})
	if err != nil {
		t.Fatalf("cpioExtract failed: %v", err)
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

func TestCpioExtract_SizeLimit(t *testing.T) {
	a := newTestArchive(t)
	a.maxSize = 20
	_, err := a.cpioExtract(filepath.Join(a.Workdir, "test.cpio"), []string{"foo/baar.txt"})
	if err == nil {
		t.Fatal("expected error for large file, but got nil")
	}
	if !strings.Contains(err.Error(), "is too large") {
		t.Fatalf("expected size limit error, got: %v", err)
	}
}

func TestTarGzExtract(t *testing.T) {
	a := newTestArchive(t)
	extractedFiles, err := a.tarGzExtract(filepath.Join(a.Workdir, "test.tar.gz"), []string{"foo/baar.txt"})
	if err != nil {
		t.Fatalf("tarGzExtract failed: %v", err)
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

func TestTarGzExtract_SizeLimit(t *testing.T) {
	a := newTestArchive(t)
	a.maxSize = 20
	_, err := a.tarGzExtract(filepath.Join(a.Workdir, "test.tar.gz"), []string{"foo/baar.txt"})
	if err == nil {
		t.Fatal("expected error for large file, but got nil")
	}
	if !strings.Contains(err.Error(), "is too large") {
		t.Fatalf("expected size limit error, got: %v", err)
	}
}

func TestTarBz2Extract(t *testing.T) {
	a := newTestArchive(t)
	extractedFiles, err := a.tarBz2Extract(filepath.Join(a.Workdir, "test.tar.bz2"), []string{"foo/baar.txt"})
	if err != nil {
		t.Fatalf("tarBz2Extract failed: %v", err)
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

func TestTarBz2Extract_SizeLimit(t *testing.T) {
	a := newTestArchive(t)
	a.maxSize = 20
	_, err := a.tarBz2Extract(filepath.Join(a.Workdir, "test.tar.bz2"), []string{"foo/baar.txt"})
	if err == nil {
		t.Fatal("expected error for large file, but got nil")
	}
	if !strings.Contains(err.Error(), "is too large") {
		t.Fatalf("expected size limit error, got: %v", err)
	}
}

func TestTarXzExtract(t *testing.T) {
	a := newTestArchive(t)
	extractedFiles, err := a.tarXzExtract(filepath.Join(a.Workdir, "test.tar.xz"), []string{"foo/baar.txt"})
	if err != nil {
		t.Fatalf("tarXzExtract failed: %v", err)
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

func TestTarXzExtract_SizeLimit(t *testing.T) {
	a := newTestArchive(t)
	a.maxSize = 20
	_, err := a.tarXzExtract(filepath.Join(a.Workdir, "test.tar.xz"), []string{"foo/baar.txt"})
	if err == nil {
		t.Fatal("expected error for large file, but got nil")
	}
	if !strings.Contains(err.Error(), "is too large") {
		t.Fatalf("expected size limit error, got: %v", err)
	}
}

func TestZipExtract(t *testing.T) {
	a := newTestArchive(t)
	extractedFiles, err := a.zipExtract(filepath.Join(a.Workdir, "test.zip"), []string{"foo/baar.txt"})
	if err != nil {
		t.Fatalf("zipExtract failed: %v", err)
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

func TestZipExtract_SizeLimit(t *testing.T) {
	a := newTestArchive(t)
	a.maxSize = 20
	_, err := a.zipExtract(filepath.Join(a.Workdir, "test.zip"), []string{"foo/baar.txt"})
	if err == nil {
		t.Fatal("expected error for large file, but got nil")
	}
	if !strings.Contains(err.Error(), "is too large") {
		t.Fatalf("expected size limit error, got: %v", err)
	}
}

func TestCpioList_Depth(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.cpioList(filepath.Join(a.Workdir, "test.cpio"), 1)
	if err != nil {
		t.Fatalf("cpioList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo", size: 0},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestTarGzList_Depth(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.tarGzList(filepath.Join(a.Workdir, "test.tar.gz"), 1)
	if err != nil {
		t.Fatalf("tarGzList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestTarBz2List_Depth(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.tarBz2List(filepath.Join(a.Workdir, "test.tar.bz2"), 1)
	if err != nil {
		t.Fatalf("tarBz2List failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestTarXzList_Depth(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.tarXzList(filepath.Join(a.Workdir, "test.tar.xz"), 1)
	if err != nil {
		t.Fatalf("tarXzList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestZipList_Depth(t *testing.T) {
	a := newTestArchive(t)
	files, err := a.zipList(filepath.Join(a.Workdir, "test.zip"), 1)
	if err != nil {
		t.Fatalf("zipList failed: %v", err)
	}

	expected := []expectedFile{
		{name: "foo/", size: 0},
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		if !containsFile(files, exp) {
			t.Errorf("expected file '%v' not found in archive", exp)
		}
	}
}

func TestSecurePath(t *testing.T) {
	a := newTestArchive(t)
	path, err := a.securePath(filepath.Join(a.Workdir, "test.zip"))
	if err != nil {
		t.Fatalf("securePath failed: %v", err)
	}
	expected, _ := filepath.Abs("../testdata/test.zip")
	if path != expected {
		t.Errorf("expected path %s, got %s", expected, path)
	}
}

func TestSecurePath_Traversal(t *testing.T) {
	a := newTestArchive(t)
	_, err := a.securePath(filepath.Join(a.Workdir, "../archive/archive.go"))
	if err == nil {
		t.Fatal("expected error for path traversal, but got nil")
	}
	if !strings.Contains(err.Error(), "is outside of the working directory") {
		t.Fatalf("expected path traversal error, got: %v", err)
	}
}

func TestSecurePath_Symlink(t *testing.T) {
	// Create a symlink from testdata/symlink to ../archive/archive.go
	// and make sure it is detected.
	a := newTestArchive(t)
	symlink := filepath.Join(a.Workdir, "symlink")
	target := "../archive/archive.go"
	err := os.Symlink(target, symlink)
	if err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}
	defer os.Remove(symlink)

	_, err = a.securePath(filepath.Join(a.Workdir, "symlink"))
	if err == nil {
		t.Fatal("expected error for symlink traversal, but got nil")
	}
	if !strings.Contains(err.Error(), "is outside of the working directory") {
		t.Fatalf("expected path traversal error, got: %v", err)
	}
}

func TestListArchiveFilesAPI(t *testing.T) {
	a := newTestArchive(t)
	archiveTypes := []string{
		"test.cpio", "test.tar.gz", "test.tar.bz2", "test.tar.xz", "test.zip",
	}

	for _, archiveType := range archiveTypes {
		t.Run(archiveType, func(t *testing.T) {
			args := ListArchiveFilesArgs{
				Path:  filepath.Join(a.Workdir, archiveType),
				Depth: 0,
			}
			session := &mcp.ServerSession{}
			_, result, err := a.ListArchiveFiles(context.Background(), &mcp.CallToolRequest{Session: session}, args)
			if err != nil {
				t.Fatalf("ListArchiveFiles failed for %s: %v", archiveType, err)
			}

			listResult, ok := result.(ListArchiveFilesResult)
			if !ok {
				t.Fatalf("unexpected result type: %T", result)
			}

			if listResult.TotalFiles < 3 {
				t.Errorf("expected at least 3 files, got %d", listResult.TotalFiles)
			}
		})
	}
}

func TestExtractArchiveFilesAPI(t *testing.T) {
	a := newTestArchive(t)
	archiveTypes := []string{
		"test.cpio", "test.tar.gz", "test.tar.bz2", "test.tar.xz", "test.zip",
	}

	for _, archiveType := range archiveTypes {
		t.Run(archiveType, func(t *testing.T) {
			args := ExtractArchiveFilesArgs{
				Path:  filepath.Join(a.Workdir, archiveType),
				Files: []string{"foo/baar.txt"},
			}
			session := &mcp.ServerSession{}
			_, result, err := a.ExtractArchiveFiles(context.Background(), &mcp.CallToolRequest{Session: session}, args)
			if err != nil {
				t.Fatalf("ExtractArchiveFiles failed for %s: %v", archiveType, err)
			}

			extractResult, ok := result.(ExtractArchiveFilesResult)
			if !ok {
				t.Fatalf("unexpected result type: %T", result)
			}
			if len(extractResult.Files) != 1 {
				t.Fatalf("expected 1 file, got %d", len(extractResult.Files))
			}
			file := extractResult.Files[0]
			if file.Name != "foo/baar.txt" {
				t.Errorf("unexpected file name: %s", file.Name)
			}
			if file.Content != "das Pferd isst Gurkensalat\n" {
				t.Errorf("unexpected content in extracted file: %s", file.Content)
			}
		})
	}
}
