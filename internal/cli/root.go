package cli

import "github.com/spf13/cobra"

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "mcp-client",
		Short:         "Trello MCP CLI - interact with Trello through MCP tools",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.PersistentFlags().Bool("json", false, "JSON output mode")
	cmd.AddCommand(newToolsCmd())
	cmd.AddCommand(newBoardsCmd())
	cmd.AddCommand(newListsCmd())
	cmd.AddCommand(newCardsCmd())
	cmd.AddCommand(newChecklistsCmd())
	cmd.AddCommand(newActionsCmd())
	cmd.AddCommand(newMembersCmd())
	cmd.AddCommand(newLabelsCmd())
	cmd.AddCommand(newOrganizationsCmd())
	return cmd
}

func jsonOut(cmd *cobra.Command) bool {
	v, _ := cmd.Root().PersistentFlags().GetBool("json")
	return v
}
