package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCardsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "cards", Short: "Manage cards"}
	cmd.AddCommand(newCardsListCmd())
	cmd.AddCommand(newCardsGetCmd())
	cmd.AddCommand(newCardsCreateCmd())
	cmd.AddCommand(newCardsUpdateCmd())
	cmd.AddCommand(newCardsMoveCmd())
	cmd.AddCommand(newCardsCommentCmd())
	cmd.AddCommand(newCardsArchiveCmd())
	cmd.AddCommand(newCardsUnarchiveCmd())
	cmd.AddCommand(newCardsDeleteCmd())
	cmd.AddCommand(newCardsSearchCmd())
	return cmd
}

func newCardsSearchCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "search", Short: "Search cards by text query"}
	query := cmd.Flags().String("query", "", "Search query")
	boardID := cmd.Flags().String("board-id", "", "Board ID (optional)")
	cmd.MarkFlagRequired("query")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		argsMap := map[string]any{"query": *query}
		if *boardID != "" {
			argsMap["board_id"] = *boardID
		}
		result, err := callTool("trello_search_cards", argsMap)
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardsList(result)
	}
	return cmd
}

func newCardsArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "archive", Short: "Archive a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_archive_card", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardDetail(result)
	}
	return cmd
}

func newCardsUnarchiveCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "unarchive", Short: "Unarchive a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_unarchive_card", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardDetail(result)
	}
	return cmd
}

func newCardsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "delete", Short: "Permanently delete a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_delete_card", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		fmt.Println("Card deleted")
		return nil
	}
	return cmd
}

func newCardsListCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list", Short: "List cards in a list"}
	listID := cmd.Flags().String("list-id", "", "List ID")
	cmd.MarkFlagRequired("list-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_list_cards", map[string]any{"list_id": *listID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardsList(result)
	}
	return cmd
}

func newCardsGetCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "get", Short: "Get card details"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_get_card", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardDetail(result)
	}
	return cmd
}

func newCardsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create", Short: "Create a card in a list"}
	listID := cmd.Flags().String("list-id", "", "List ID")
	name := cmd.Flags().String("name", "", "Card name")
	desc := cmd.Flags().String("desc", "", "Card description")
	due := cmd.Flags().String("due", "", "Due date (RFC3339)")
	cmd.MarkFlagRequired("list-id")
	cmd.MarkFlagRequired("name")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		argsMap := map[string]any{"list_id": *listID, "name": *name}
		if *desc != "" {
			argsMap["desc"] = *desc
		}
		if *due != "" {
			argsMap["due"] = *due
		}
		result, err := callTool("trello_create_card", argsMap)
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardDetail(result)
	}
	return cmd
}

func newCardsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "update", Short: "Update a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	name := cmd.Flags().String("name", "", "New card name")
	desc := cmd.Flags().String("desc", "", "New card description")
	due := cmd.Flags().String("due", "", "Due date (RFC3339)")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		argsMap := map[string]any{"card_id": *cardID}
		if *name != "" {
			argsMap["name"] = *name
		}
		if *desc != "" {
			argsMap["desc"] = *desc
		}
		if *due != "" {
			argsMap["due"] = *due
		}
		if len(argsMap) < 2 {
			return fmtUsage("at least one of --name, --desc, or --due must be provided")
		}
		result, err := callTool("trello_update_card", argsMap)
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardDetail(result)
	}
	return cmd
}

func newCardsMoveCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "move", Short: "Move a card to another list"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	listID := cmd.Flags().String("list-id", "", "Target list ID")
	cmd.MarkFlagRequired("card-id")
	cmd.MarkFlagRequired("list-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_move_card", map[string]any{"card_id": *cardID, "list_id": *listID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardDetail(result)
	}
	return cmd
}

func newCardsCommentCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "comment", Short: "Add a comment to a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	text := cmd.Flags().String("text", "", "Comment text")
	cmd.MarkFlagRequired("card-id")
	cmd.MarkFlagRequired("text")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_add_comment", map[string]any{"card_id": *cardID, "text": *text})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printCardDetail(result)
	}
	return cmd
}
