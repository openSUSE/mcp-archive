package archive

import (
	"context"
	"fmt"
	"os"

	"github.com/cavaliergopher/cpio"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListArchiveFilesArgs are the arguments for the list_archive_files tool.
type ListArchiveFilesArgs struct {
	Path string `json:"path" jsonschema:"the path to the cpio archive"`
}

// ListArchiveFiles lists the files in a cpio archive.
func ListArchiveFiles(ctx context.Context, req *mcp.CallToolRequest, args ListArchiveFilesArgs) (*mcp.CallToolResult, any, error) {
	file, err := os.Open(args.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	reader := cpio.NewReader(file)
	var files []string
	for {
		header, err := reader.Next()
		if err != nil {
			break
		}
		files = append(files, fmt.Sprintf("%s %d %s", header.Name, header.Size, header.Mode))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("%v", files)},
		},
	}, nil, nil
}
