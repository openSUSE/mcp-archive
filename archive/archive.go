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
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cavaliergopher/cpio"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ulikunitz/xz"
)

// Archive holds the configuration for the archive tools.
type Archive struct {
	maxSize int64
	Workdir string
}

// New creates a new Archive instance.
func New(workdir string) (*Archive, error) {
	absWorkdir, err := filepath.Abs(workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for workdir: %w", err)
	}
	return &Archive{
		maxSize: 100 * 1024,
		Workdir: absWorkdir,
	}, nil
}

// FileInfo represents a file in an archive.
type FileInfo struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Permissions string `json:"permissions"`
}

// ListArchiveFilesArgs are the arguments for the list_archive_files tool.
type ListArchiveFilesArgs struct {
	Path           string `json:"path" jsonschema:"the path to the archive"`
	Depth          int    `json:"depth" jsonschema:"the depth of the directory tree to list. 0 means the complete directory tree"`
	Limit          int    `json:"limit,omitempty" jsonschema:"the maximum number of files to display. If not set, it will default to 100"`
	IncludePattern string `json:"include,omitempty" jsonschema:"an optional regular expression to include files"`
	ExcludePattern string `json:"exclude,omitempty" jsonschema:"an optional regular expression to exclude files"`
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

func (a *Archive) securePath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("path is not an absolute path: %s", path)
	}
	absPath := filepath.Clean(path)
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate symlinks: %w", err)
	}

	if !strings.HasPrefix(evalPath, a.Workdir) {
		return "", fmt.Errorf("path %s is outside of the working directory", path)
	}
	return evalPath, nil
}

func (a *Archive) cpioList(path string, depth int) ([]FileInfo, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	reader := cpio.NewReader(file)
	var files []FileInfo
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if depth > 0 && len(strings.Split(strings.Trim(header.Name, "/"), "/")) > depth {
			continue
		}
		files = append(files, FileInfo{
			Name:        header.Name,
			Size:        header.Size,
			Permissions: header.Mode.String(),
		})
	}
	return files, nil
}

func (a *Archive) tarGzList(path string, depth int) ([]FileInfo, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
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
	var files []FileInfo
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if depth > 0 && len(strings.Split(strings.Trim(header.Name, "/"), "/")) > depth {
			continue
		}
		files = append(files, FileInfo{
			Name:        header.Name,
			Size:        header.Size,
			Permissions: os.FileMode(header.Mode).String(),
		})
	}
	return files, nil
}

func (a *Archive) tarBz2List(path string, depth int) ([]FileInfo, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	bz2r := bzip2.NewReader(file)
	tr := tar.NewReader(bz2r)
	var files []FileInfo
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if depth > 0 && len(strings.Split(strings.Trim(header.Name, "/"), "/")) > depth {
			continue
		}
		files = append(files, FileInfo{
			Name:        header.Name,
			Size:        header.Size,
			Permissions: os.FileMode(header.Mode).String(),
		})
	}
	return files, nil
}

func (a *Archive) tarXzList(path string, depth int) ([]FileInfo, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	xzr, err := xz.NewReader(file)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(xzr)
	var files []FileInfo
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if depth > 0 && len(strings.Split(strings.Trim(header.Name, "/"), "/")) > depth {
			continue
		}
		files = append(files, FileInfo{
			Name:        header.Name,
			Size:        header.Size,
			Permissions: os.FileMode(header.Mode).String(),
		})
	}
	return files, nil
}

func (a *Archive) zipList(path string, depth int) ([]FileInfo, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	r, err := zip.OpenReader(securePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var files []FileInfo
	for _, f := range r.File {
		if depth > 0 && len(strings.Split(strings.Trim(f.Name, "/"), "/")) > depth {
			continue
		}
		files = append(files, FileInfo{
			Name:        f.Name,
			Size:        int64(f.UncompressedSize64),
			Permissions: f.Mode().String(),
		})
	}
	return files, nil
}

// ListArchiveFilesResult holds the result of the list_archive_files tool.
type ListArchiveFilesResult struct {
	TotalFiles     int        `json:"total_files"`
	FilteredFiles  int        `json:"filtered_files"`
	DisplayedFiles int        `json:"displayed_files"`
	Files          []FileInfo `json:"files"`
}

// ListArchiveFiles lists the files in an archive.
func (a *Archive) ListArchiveFiles(ctx context.Context, req *mcp.CallToolRequest, args ListArchiveFilesArgs) (*mcp.CallToolResult, any, error) {
	slog.Debug("mcp tool call: ListArchiveFiles", "session", req.Session.ID(), "params", args)
	var files []FileInfo
	var err error

	switch {
	case strings.HasSuffix(args.Path, ".cpio"):
		files, err = a.cpioList(args.Path, args.Depth)
	case strings.HasSuffix(args.Path, ".tar.gz"):
		files, err = a.tarGzList(args.Path, args.Depth)
	case strings.HasSuffix(args.Path, ".tar.bz2"):
		files, err = a.tarBz2List(args.Path, args.Depth)
	case strings.HasSuffix(args.Path, ".tar.xz"):
		files, err = a.tarXzList(args.Path, args.Depth)
	case strings.HasSuffix(args.Path, ".zip"):
		files, err = a.zipList(args.Path, args.Depth)
	default:
		return nil, nil, fmt.Errorf("unsupported archive format for %s", args.Path)
	}

	if err != nil {
		return nil, nil, err
	}

	totalFiles := len(files)
	var filteredFiles []FileInfo

	for _, file := range files {
		includeMatch := true
		if args.IncludePattern != "" {
			includeMatch, err = regexp.MatchString(args.IncludePattern, file.Name)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid include pattern: %w", err)
			}
		}

		excludeMatch := false
		if args.ExcludePattern != "" {
			excludeMatch, err = regexp.MatchString(args.ExcludePattern, file.Name)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid exclude pattern: %w", err)
			}
		}

		if includeMatch && !excludeMatch {
			filteredFiles = append(filteredFiles, file)
		}
	}

	limit := args.Limit
	if limit == 0 {
		limit = 100
	}

	displayedFilesCount := len(filteredFiles)
	if displayedFilesCount > limit {
		displayedFilesCount = limit
	}

	result := ListArchiveFilesResult{
		TotalFiles:     totalFiles,
		FilteredFiles:  len(filteredFiles),
		DisplayedFiles: displayedFilesCount,
		Files:          filteredFiles[:displayedFilesCount],
	}

	return nil, result, nil
}

func (a *Archive) cpioExtract(path string, filesToExtract []string) ([]File, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
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
				if header.Size > a.maxSize {
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

func (a *Archive) tarGzExtract(path string, filesToExtract []string) ([]File, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
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
				if header.Size > a.maxSize {
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

func (a *Archive) tarBz2Extract(path string, filesToExtract []string) ([]File, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
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
				if header.Size > a.maxSize {
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

func (a *Archive) tarXzExtract(path string, filesToExtract []string) ([]File, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(securePath)
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
				if header.Size > a.maxSize {
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

func (a *Archive) zipExtract(path string, filesToExtract []string) ([]File, error) {
	securePath, err := a.securePath(path)
	if err != nil {
		return nil, err
	}
	r, err := zip.OpenReader(securePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var extractedFiles []File
	for _, f := range r.File {
		for _, fileToExtract := range filesToExtract {
			if f.Name == fileToExtract {
				if f.UncompressedSize64 > uint64(a.maxSize) {
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

// ExtractArchiveFilesResult holds the result of the extract_archive_files tool.
type ExtractArchiveFilesResult struct {
	Files []File `json:"files"`
}

// ExtractArchiveFiles extracts files from an archive and returns their content.
func (a *Archive) ExtractArchiveFiles(ctx context.Context, req *mcp.CallToolRequest, args ExtractArchiveFilesArgs) (*mcp.CallToolResult, any, error) {
	slog.Debug("mcp tool call: ExtractArchiveFiles", "session", req.Session.ID(), "params", args)
	var files []File
	var err error

	switch {
	case strings.HasSuffix(args.Path, ".cpio"):
		files, err = a.cpioExtract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".tar.gz"):
		files, err = a.tarGzExtract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".tar.bz2"):
		files, err = a.tarBz2Extract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".tar.xz"):
		files, err = a.tarXzExtract(args.Path, args.Files)
	case strings.HasSuffix(args.Path, ".zip"):
		files, err = a.zipExtract(args.Path, args.Files)
	default:
		return nil, nil, fmt.Errorf("unsupported archive format for %s", args.Path)
	}

	if err != nil {
		return nil, nil, err
	}

	return nil, ExtractArchiveFilesResult{Files: files}, nil
}
