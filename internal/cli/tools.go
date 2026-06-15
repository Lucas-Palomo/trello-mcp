package cli

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

func newToolsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "tools", Short: "Manage MCP tools"}
	cmd.AddCommand(&cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := mcpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
			if err != nil {
				return err
			}
			if jsonOut(cmd) {
				return printJSON(result)
			}
			return printToolsList(result)
		},
	})
	return cmd
}
