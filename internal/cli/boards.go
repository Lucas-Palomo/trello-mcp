package cli

import (
	"github.com/spf13/cobra"
)

func newBoardsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "boards", Short: "Manage boards"}
	cmd.AddCommand(newBoardsListCmd())
	cmd.AddCommand(newBoardsGetCmd())
	cmd.AddCommand(newBoardsCreateCmd())
	return cmd
}

func newBoardsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create", Short: "Create a new board"}
	name := cmd.Flags().String("name", "", "Board name")
	noLists := cmd.Flags().Bool("no-default-lists", false, "Do not create default lists")
	cmd.MarkFlagRequired("name")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		argsMap := map[string]any{"name": *name}
		if *noLists {
			argsMap["default_lists"] = false
		}
		result, err := callTool("trello_create_board", argsMap)
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printBoardDetail(result)
	}
	return cmd
}

func newBoardsListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := callTool("trello_list_boards", nil)
			if err != nil {
				return err
			}
			if result.IsError {
				return fmtToolError(result)
			}
			if jsonOut(cmd) {
				return printJSON(result)
			}
			return printBoardsList(result)
		},
	}
}

func newBoardsGetCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "get", Short: "Get board details"}
	flags := cmd.Flags()
	boardID := flags.String("board-id", "", "Board ID")
	cmd.MarkFlagRequired("board-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_get_board", map[string]any{"board_id": *boardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printBoardDetail(result)
	}
	return cmd
}
