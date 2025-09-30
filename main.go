package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openSUSE/mcp-archive/internal/archive"
)

var (
	httpAddr = flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
)

func main() {
	// Create a server with a single tool that says "Hi".
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter"}, nil)

	// Add the tools from the hello package.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_archive_files",
		Description: "list the files in a cpio archive",
	}, archive.ListArchiveFiles)

	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP handler listening at %s", *httpAddr)
		log.Fatal(http.ListenAndServe(*httpAddr, handler))
	} else {
		t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
		if err := server.Run(context.Background(), t); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}
}
