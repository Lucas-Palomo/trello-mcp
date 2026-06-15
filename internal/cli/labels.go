package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLabelsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "labels", Short: "Manage labels"}
	cmd.AddCommand(newLabelsBoardListCmd())
	cmd.AddCommand(newLabelsCardListCmd())
	cmd.AddCommand(newLabelsAddCmd())
	cmd.AddCommand(newLabelsRemoveCmd())
	return cmd
}

func newLabelsAddCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "add", Short: "Add a label to a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	labelID := cmd.Flags().String("label-id", "", "Label ID")
	cmd.MarkFlagRequired("card-id")
	cmd.MarkFlagRequired("label-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_add_card_label", map[string]any{"card_id": *cardID, "label_id": *labelID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		fmt.Println("Label added")
		return nil
	}
	return cmd
}

func newLabelsRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "remove", Short: "Remove a label from a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	labelID := cmd.Flags().String("label-id", "", "Label ID")
	cmd.MarkFlagRequired("card-id")
	cmd.MarkFlagRequired("label-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_remove_card_label", map[string]any{"card_id": *cardID, "label_id": *labelID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		fmt.Println("Label removed")
		return nil
	}
	return cmd
}

func newLabelsBoardListCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "board-list", Short: "List labels on a board"}
	boardID := cmd.Flags().String("board-id", "", "Board ID")
	cmd.MarkFlagRequired("board-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_list_board_labels", map[string]any{"board_id": *boardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printLabelsList(result)
	}
	return cmd
}

func newLabelsCardListCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "card-list", Short: "List labels on a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_list_card_labels", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardLabelsList(result)
	}
	return cmd
}
