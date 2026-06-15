package cli

import (
	"github.com/spf13/cobra"
)

func newChecklistsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "checklists", Short: "Manage checklists"}
	cmd.AddCommand(newChecklistsListCmd())
	return cmd
}

func newChecklistsListCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list", Short: "List checklists on a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_list_checklists", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printChecklistsList(result)
	}
	return cmd
}
