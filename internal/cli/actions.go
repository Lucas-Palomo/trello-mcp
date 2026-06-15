package cli

import (
	"github.com/spf13/cobra"
)

func newActionsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "actions", Short: "Manage card actions"}
	cmd.AddCommand(newActionsListCmd())
	return cmd
}

func newActionsListCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list", Short: "List actions on a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_list_card_actions", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printActionsList(result)
	}
	return cmd
}
