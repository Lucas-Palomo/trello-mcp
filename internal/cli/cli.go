package cli

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func fmtToolError(result *mcp.CallToolResult) error {
	return fmt.Errorf("tool error: %s", extractError(result))
}

func fmtUsage(msg string) error {
	return fmt.Errorf("usage: %s", msg)
}

var mcpClient *client.Client

func Run(args []string) error {
	c, err := newMCPClient()
	if err != nil {
		return fmt.Errorf("server setup failed: %w", err)
	}
	defer c.Close()
	mcpClient = c

	root := newRootCmd()
	root.SetArgs(args)
	return root.Execute()
}

func runWithClient(args []string, c *client.Client) error {
	mcpClient = c
	root := newRootCmd()
	root.SetArgs(args)
	return root.Execute()
}

func callTool(toolName string, args map[string]any) (*mcp.CallToolResult, error) {
	return mcpClient.CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{Name: toolName, Arguments: args},
	})
}

func extractError(result *mcp.CallToolResult) string {
	for _, c := range result.Content {
		if text, ok := c.(mcp.TextContent); ok {
			return text.Text
		}
	}
	return "unknown error"
}
