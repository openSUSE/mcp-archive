// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package archive

import (
	"reflect"
	"testing"
)

func TestCpioList(t *testing.T) {
	files, err := cpioList("../../testdata/test.cpio")
	if err != nil {
		t.Fatalf("cpioList failed: %v", err)
	}

	expected := []string{
		"foo 0 040755",
		"foo/baar.txt 5 0100644",
		"foo/bazz 5 0100644",
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
			if reflect.DeepEqual(file, exp) {
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
	files, err := tarGzList("../../testdata/test.tar.gz")
	if err != nil {
		t.Fatalf("tarGzList failed: %v", err)
	}

	expected := []string{
		"foo/ 0 -rwxr-xr-x",
		"foo/baar.txt 5 -rw-r--r--",
		"foo/bazz 5 -rw-r--r--",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if reflect.DeepEqual(file, exp) {
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
	files, err := tarBz2List("../../testdata/test.tar.bz2")
	if err != nil {
		t.Fatalf("tarBz2List failed: %v", err)
	}

	expected := []string{
		"foo/ 0 -rwxr-xr-x",
		"foo/baar.txt 5 -rw-r--r--",
		"foo/bazz 5 -rw-r--r--",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if reflect.DeepEqual(file, exp) {
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
	files, err := tarXzList("../../testdata/test.tar.xz")
	if err != nil {
		t.Fatalf("tarXzList failed: %v", err)
	}

	expected := []string{
		"foo/ 0 -rwxr-xr-x",
		"foo/baar.txt 5 -rw-r--r--",
		"foo/bazz 5 -rw-r--r--",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if reflect.DeepEqual(file, exp) {
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
	files, err := zipList("../../testdata/test.zip")
	if err != nil {
		t.Fatalf("zipList failed: %v", err)
	}

	expected := []string{
		"foo/ 0 drwxr-xr-x",
		"foo/baar.txt 5 -rw-r--r--",
		"foo/bazz 5 -rw-r--r--",
	}

	// The order of files in the archive is not guaranteed, so we need to compare them in a way that ignores order.
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, exp := range expected {
		found := false
		for _, file := range files {
			if reflect.DeepEqual(file, exp) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file '%s' not found in archive", exp)
		}
	}
}
