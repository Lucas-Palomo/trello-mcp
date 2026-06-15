package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newMembersCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "members", Short: "Manage card members"}
	cmd.AddCommand(newMembersListCmd())
	cmd.AddCommand(newMembersAddCmd())
	cmd.AddCommand(newMembersRemoveCmd())
	return cmd
}

func newMembersAddCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "add", Short: "Assign a member to a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	memberID := cmd.Flags().String("member-id", "", "Member ID")
	cmd.MarkFlagRequired("card-id")
	cmd.MarkFlagRequired("member-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_add_card_member", map[string]any{"card_id": *cardID, "member_id": *memberID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		fmt.Println("Member added")
		return nil
	}
	return cmd
}

func newMembersRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "remove", Short: "Remove a member from a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	memberID := cmd.Flags().String("member-id", "", "Member ID")
	cmd.MarkFlagRequired("card-id")
	cmd.MarkFlagRequired("member-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_remove_card_member", map[string]any{"card_id": *cardID, "member_id": *memberID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		fmt.Println("Member removed")
		return nil
	}
	return cmd
}

func newMembersListCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list", Short: "List members on a card"}
	cardID := cmd.Flags().String("card-id", "", "Card ID")
	cmd.MarkFlagRequired("card-id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := callTool("trello_list_card_members", map[string]any{"card_id": *cardID})
		if err != nil {
			return err
		}
		if result.IsError {
			return fmtToolError(result)
		}
		if jsonOut(cmd) {
			return printJSON(result)
		}
		return printMembersList(result)
	}
	return cmd
}
