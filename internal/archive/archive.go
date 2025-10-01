package archive

import (
	"archive/tar"
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
