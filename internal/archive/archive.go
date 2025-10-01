// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cavaliergopher/cpio"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ulikunitz/xz"
)

// ListArchiveFilesArgs are the arguments for the list_archive_files tool.
type ListArchiveFilesArgs struct {
	Path string `json:"path" jsonschema:"the path to the archive"`
}

// ExtractArchiveFilesArgs are the arguments for the extract_archive_files tool.
type ExtractArchiveFilesArgs struct {
	Path  string   `json:"path" jsonschema:"the path to the archive"`
	Files []string `json:"files" jsonschema:"the files to extract"`
}

// File represents an extracted file's content and metadata.
type File struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Permissions string `json:"permissions"`
	Content     string `json:"content"`
}

var (
	// MaxExtractFileSize is the maximum size of a file that can be extracted.
	MaxExtractFileSize int64 = 100 * 1024
)

func cpioList(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	reader := cpio.NewReader(file)
	var files []string
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, fmt.Sprintf("%s %d %s", header.Name, header.Size, header.Mode))
	}
	return files, nil
}

func tarGzList(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var files []string
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, fmt.Sprintf("%s %d %s", header.Name, header.Size, os.FileMode(header.Mode).String()))
	}
	return files, nil
}

func tarBz2List(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	bz2r := bzip2.NewReader(file)
	tr := tar.NewReader(bz2r)
	var files []string
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, fmt.Sprintf("%s %d %s", header.Name, header.Size, os.FileMode(header.Mode).String()))
	}
	return files, nil
}

func tarXzList(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	xzr, err := xz.NewReader(file)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(xzr)
	var files []string
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, fmt.Sprintf("%s %d %s", header.Name, header.Size, os.FileMode(header.Mode).String()))
	}
	return files, nil
}

func zipList(path string) ([]string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var files []string
	for _, f := range r.File {
		files = append(files, fmt.Sprintf("%s %d %s", f.Name, f.UncompressedSize64, f.Mode().String()))
	}
	return files, nil
}

// ListArchiveFiles lists the files in an archive.
func ListArchiveFiles(ctx context.Context, req *mcp.CallToolRequest, args ListArchiveFilesArgs) (*mcp.CallToolResult, any, error) {
	var files []string
	var err error

	switch {
	case strings.HasSuffix(args.Path, ".cpio"):
		files, err = cpioList(args.Path)
	case strings.HasSuffix(args.Path, ".tar.gz"):
		files, err = tarGzList(args.Path)
	case strings.HasSuffix(args.Path, ".tar.bz2"):
		files, err = tarBz2List(args.Path)
	case strings.HasSuffix(args.Path, ".tar.xz"):
		files, err = tarXzList(args.Path)
	case strings.HasSuffix(args.Path, ".zip"):
		files, err = zipList(args.Path)
	default:
		return nil, nil, fmt.Errorf("unsupported archive format for %s", args.Path)
	}

	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("%v", files)},
		},
	}, nil, nil
}

func cpioExtract(path string, filesToExtract []string) ([]File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	reader := cpio.NewReader(file)
	var extractedFiles []File

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, f := range filesToExtract {
			if header.Name == f {
				if header.Size > MaxExtractFileSize {
					return nil, fmt.Errorf("file %s is too large to extract: %d bytes", header.Name, header.Size)
				}

				buf := make([]byte, header.Size)
				if _, err := io.ReadFull(reader, buf); err != nil {
					return nil, fmt.Errorf("could not read file %s from archive: %w", header.Name, err)
				}

				extractedFile := File{
					Name:        header.Name,
					Size:        header.Size,
					Permissions: header.Mode.String(),
					Content:     string(buf),
				}
				extractedFiles = append(extractedFiles, extractedFile)
			}
		}
	}
	return extractedFiles, nil
}

func tarGzExtract(path string, filesToExtract []string) ([]File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var extractedFiles []File

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, f := range filesToExtract {
			if header.Name == f {
				if header.Size > MaxExtractFileSize {
					return nil, fmt.Errorf("file %s is too large to extract: %d bytes", header.Name, header.Size)
				}

				buf := make([]byte, header.Size)
				if _, err := io.ReadFull(tr, buf); err != nil {
					return nil, fmt.Errorf("could not read file %s from archive: %w", header.Name, err)
				}

				extractedFile := File{
					Name:        header.Name,
					Size:        header.Size,
					Permissions: os.FileMode(header.Mode).String(),
					Content:     string(buf),
				}
				extractedFiles = append(extractedFiles, extractedFile)
			}
		}
	}
	return extractedFiles, nil
}

func tarBz2Extract(path string, filesToExtract []string) ([]File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	bz2r := bzip2.NewReader(file)
	tr := tar.NewReader(bz2r)
	var extractedFiles []File

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, f := range filesToExtract {
			if header.Name == f {
				if header.Size > MaxExtractFileSize {
					return nil, fmt.Errorf("file %s is too large to extract: %d bytes", header.Name, header.Size)
				}

				buf := make([]byte, header.Size)
				if _, err := io.ReadFull(tr, buf); err != nil {
					return nil, fmt.Errorf("could not read file %s from archive: %w", header.Name, err)
				}

				extractedFile := File{
					Name:        header.Name,
					Size:        header.Size,
					Permissions: os.FileMode(header.Mode).String(),
					Content:     string(buf),
				}
				extractedFiles = append(extractedFiles, extractedFile)
			}
		}
	}
	return extractedFiles, nil
}

func tarXzExtract(path string, filesToExtract []string) ([]File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	xzr, err := xz.NewReader(file)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(xzr)
	var extractedFiles []File

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, f := range filesToExtract {
			if header.Name == f {
				if header.Size > MaxExtractFileSize {
					return nil, fmt.Errorf("file %s is too large to extract: %d bytes", header.Name, header.Size)
				}

				buf := make([]byte, header.Size)
				if _, err := io.ReadFull(tr, buf); err != nil {
					return nil, fmt.Errorf("could not read file %s from archive: %w", header.Name, err)
				}

				extractedFile := File{
					Name:        header.Name,
					Size:        header.Size,
					Permissions: os.FileMode(header.Mode).String(),
					Content:     string(buf),
				}
				extractedFiles = append(extractedFiles, extractedFile)
			}
		}
	}
	return extractedFiles, nil
}

func zipExtract(path string, filesToExtract []string) ([]File, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var extractedFiles []File
	for _, f := range r.File {
		for _, fileToExtract := range filesToExtract {
			if f.Name == fileToExtract {
				if f.UncompressedSize64 > uint64(MaxExtractFileSize) {
					return nil, fmt.Errorf("file %s is too large to extract: %d bytes", f.Name, f.UncompressedSize64)
				}

				rc, err := f.Open()
				if err != nil {
					return nil, err
				}

				buf := make([]byte, f.UncompressedSize64)
				if _, err := io.ReadFull(rc, buf); err != nil {
					rc.Close()
					return nil, fmt.Errorf("could not read file %s from archive: %w", f.Name, err)
				}
				rc.Close()

				extractedFile := File{
					Name:        f.Name,
					Size:        int64(f.UncompressedSize64),
					Permissions: f.Mode().String(),
					Content:     string(buf),
				}
				extractedFiles = append(extractedFiles, extractedFile)
			}
		}
	}
	return extractedFiles, nil
}

// ExtractArchiveFiles extracts files from an archive and returns their content.
func ExtractArchiveFiles(ctx context.Context, req *mcp.CallToolRequest, args ExtractArchiveFilesArgs) (*mcp.CallToolResult, any, error) {
	var files []File
	var err error

	switch {
	case strings.HasSuffix(args.Path, ".cpio"):
		files, err = cpioExtract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".tar.gz"):
		files, err = tarGzExtract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".tar.bz2"):
		files, err = tarBz2Extract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".tar.xz"):
		files, err = tarXzExtract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".zip"):
		files, err = zipExtract(args.Path, args.Files)
	default:
		return nil, nil, fmt.Errorf("unsupported archive format for %s", args.Path)
	}

	if err != nil {
		return nil, nil, err
	}

	return nil, files, nil
}
