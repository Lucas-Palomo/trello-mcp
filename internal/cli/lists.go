package cli

import (
	"github.com/spf13/cobra"
)

func newListsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "lists", Short: "Manage lists"}
	cmd.AddCommand(newListsListCmd())
	cmd.AddCommand(newListsCreateCmd())
	return cmd
}

func newListsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create", Short: "Create a new list in a board"}
	name := cmd.Flags().String("name", "", "List name")
	boardID := cmd.Flags().String("board-id", "", "Board ID")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("board-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_create_list", map[string]any{"name": *name, "board_id": *boardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printListDetail(result)
	}
	return cmd
}

func newListsListCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list", Short: "List lists in a board"}
	boardID := cmd.Flags().String("board-id", "", "Board ID")
	cmd.MarkFlagRequired("board-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_list_lists", map[string]any{"board_id": *boardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printListsList(result)
	}
	return cmd
}
