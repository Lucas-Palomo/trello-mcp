package cli

import "github.com/spf13/cobra"

func newOrganizationsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "orgs", Short: "Manage organizations"}
	cmd.AddCommand(newOrganizationsListCmd())
	return cmd
}

func newOrganizationsListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := callTool("trello_list_organizations", nil)
			if err != nil {
				return err
			}
			if result.IsError {
				return fmtToolError(result)
			}
			if jsonOut(cmd) {
				return printJSON(result)
			}
			return printOrganizationsList(result)
		},
	}
}
