package cli

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func testMCPClient(t *testing.T) *client.Client {
	t.Helper()

	s := server.NewMCPServer("test", "0.1.0")

	addTool := func(name string, handler func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
		s.AddTool(mcp.NewTool(name), handler)
	}

	addTool("trello_list_boards", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"boards": []map[string]any{{"id": "b1", "name": "Board 1", "url": "https://trello.com/b/b1"}},
		})
	})

	addTool("trello_get_board", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"board": map[string]any{"id": "b1", "name": "Board 1", "url": "https://trello.com/b/b1"},
		})
	})

	addTool("trello_list_lists", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"lists": []map[string]any{{"id": "l1", "name": "List 1"}},
		})
	})

	addTool("trello_list_cards", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"cards": []map[string]any{{"id": "c1", "name": "Card 1", "desc": "A card"}},
		})
	})

	addTool("trello_get_card", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"card": map[string]any{"id": "c1", "name": "Card 1", "desc": "A card"},
		})
	})

	addTool("trello_create_card", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"card": map[string]any{"id": "c1", "name": "New Card"},
		})
	})

	addTool("trello_update_card", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"card": map[string]any{"id": "c1", "name": "Updated"},
		})
	})

	addTool("trello_move_card", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"card": map[string]any{"id": "c1", "name": "Moved"},
		})
	})

	addTool("trello_add_comment", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"action": map[string]any{"id": "a1", "type": "commentCard", "text": "Nice!"},
		})
	})

	addTool("trello_list_checklists", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"checklists": []map[string]any{{"id": "chk1", "name": "Checklist 1"}},
		})
	})

	addTool("trello_list_card_actions", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"actions": []map[string]any{{"id": "a1", "type": "commentCard", "text": "Nice!"}},
		})
	})

	addTool("trello_list_board_labels", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"labels": []map[string]any{{"id": "lb1", "name": "urgent", "color": "red"}},
		})
	})

	addTool("trello_list_card_labels", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"labels": []map[string]any{{"id": "lb1", "name": "bug", "color": "orange"}},
		})
	})

	addTool("trello_add_card_label", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{"success": true})
	})

	addTool("trello_remove_card_label", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{"success": true})
	})

	addTool("trello_add_card_member", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{"success": true})
	})

	addTool("trello_remove_card_member", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{"success": true})
	})

	addTool("trello_create_board", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"board": map[string]any{"id": "b1", "name": "My Board", "url": "https://trello.com/b/b1"},
		})
	})

	addTool("trello_create_list", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"list": map[string]any{"id": "l1", "name": "To Do"},
		})
	})

	addTool("trello_archive_card", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"card": map[string]any{"id": "c1", "closed": true},
		})
	})

	addTool("trello_unarchive_card", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"card": map[string]any{"id": "c1", "closed": false},
		})
	})

	addTool("trello_delete_card", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{"success": true})
	})

	addTool("trello_search_cards", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"cards": []map[string]any{{"id": "c1", "name": "Found card"}},
		})
	})

	addTool("trello_list_organizations", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"organizations": []map[string]any{{"id": "o1", "name": "acmecorp", "displayName": "Acme Corp", "url": "https://trello.com/acmecorp"}},
		})
	})

	addTool("trello_list_card_members", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultJSON(map[string]any{
			"members": []map[string]any{{"id": "m1", "full_name": "Jane Doe", "username": "jane"}},
		})
	})

	c, err := client.NewInProcessClient(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })

	_, err = c.Initialize(context.Background(), mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{Name: "test", Version: "0.1.0"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	return c
}

// Basic

func TestRun_MissingArgs(t *testing.T) {
	err := runWithClient([]string{}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	err := runWithClient([]string{"unknown"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_Help(t *testing.T) {
	err := runWithClient([]string{"--help"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// Boards

func TestRun_BoardsList_Success(t *testing.T) {
	err := runWithClient([]string{"boards", "list"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_BoardsList_JSON(t *testing.T) {
	err := runWithClient([]string{"boards", "list", "--json"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_BoardsGet_Success(t *testing.T) {
	err := runWithClient([]string{"boards", "get", "--board-id", "b1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_BoardsGet_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"boards", "get"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Lists

func TestRun_ListsList_Success(t *testing.T) {
	err := runWithClient([]string{"lists", "list", "--board-id", "b1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_ListsList_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"lists", "list"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Cards

func TestRun_CardsList_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "list", "--list-id", "l1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsList_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"cards", "list"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsList_JSON(t *testing.T) {
	err := runWithClient([]string{"cards", "list", "--list-id", "l1", "--json"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsGet_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "get", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsGet_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"cards", "get"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsCreate_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "create", "--list-id", "l1", "--name", "New Card"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsCreate_MissingListID(t *testing.T) {
	err := runWithClient([]string{"cards", "create", "--name", "Card"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsCreate_MissingName(t *testing.T) {
	err := runWithClient([]string{"cards", "create", "--list-id", "l1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsUpdate_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "update", "--card-id", "c1", "--name", "Updated"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsUpdate_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"cards", "update", "--name", "x"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsMove_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "move", "--card-id", "c1", "--list-id", "l2"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsMove_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"cards", "move", "--list-id", "l2"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsMove_MissingListID(t *testing.T) {
	err := runWithClient([]string{"cards", "move", "--card-id", "c1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsComment_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "comment", "--card-id", "c1", "--text", "Nice!"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsComment_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"cards", "comment", "--text", "hi"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsComment_MissingText(t *testing.T) {
	err := runWithClient([]string{"cards", "comment", "--card-id", "c1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Checklists

func TestRun_ChecklistsList_Success(t *testing.T) {
	err := runWithClient([]string{"checklists", "list", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_ChecklistsList_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"checklists", "list"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Actions

func TestRun_ActionsList_Success(t *testing.T) {
	err := runWithClient([]string{"actions", "list", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_ActionsList_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"actions", "list"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Members

func TestRun_MembersList_Success(t *testing.T) {
	err := runWithClient([]string{"members", "list", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_MembersList_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"members", "list"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Tools

func TestRun_ToolsList_Success(t *testing.T) {
	err := runWithClient([]string{"tools", "list"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_ToolsList_JSON(t *testing.T) {
	err := runWithClient([]string{"tools", "list", "--json"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// Labels

func TestRun_LabelsBoardList_Success(t *testing.T) {
	err := runWithClient([]string{"labels", "board-list", "--board-id", "b1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_LabelsBoardList_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"labels", "board-list"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_LabelsCardList_Success(t *testing.T) {
	err := runWithClient([]string{"labels", "card-list", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_LabelsCardList_MissingFlag(t *testing.T) {
	err := runWithClient([]string{"labels", "card-list"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_LabelsBoardList_JSON(t *testing.T) {
	err := runWithClient([]string{"labels", "board-list", "--board-id", "b1", "--json"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_LabelsAdd_Success(t *testing.T) {
	err := runWithClient([]string{"labels", "add", "--card-id", "c1", "--label-id", "lb1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_LabelsAdd_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"labels", "add", "--label-id", "lb1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_LabelsAdd_MissingLabelID(t *testing.T) {
	err := runWithClient([]string{"labels", "add", "--card-id", "c1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_LabelsRemove_Success(t *testing.T) {
	err := runWithClient([]string{"labels", "remove", "--card-id", "c1", "--label-id", "lb1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_LabelsRemove_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"labels", "remove", "--label-id", "lb1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_LabelsRemove_MissingLabelID(t *testing.T) {
	err := runWithClient([]string{"labels", "remove", "--card-id", "c1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Members mutations

func TestRun_MembersAdd_Success(t *testing.T) {
	err := runWithClient([]string{"members", "add", "--card-id", "c1", "--member-id", "m1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_MembersAdd_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"members", "add", "--member-id", "m1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_MembersAdd_MissingMemberID(t *testing.T) {
	err := runWithClient([]string{"members", "add", "--card-id", "c1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_MembersRemove_Success(t *testing.T) {
	err := runWithClient([]string{"members", "remove", "--card-id", "c1", "--member-id", "m1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_MembersRemove_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"members", "remove", "--member-id", "m1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_MembersRemove_MissingMemberID(t *testing.T) {
	err := runWithClient([]string{"members", "remove", "--card-id", "c1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Boards create

func TestRun_BoardsCreate_Success(t *testing.T) {
	err := runWithClient([]string{"boards", "create", "--name", "My Board"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_BoardsCreate_MissingName(t *testing.T) {
	err := runWithClient([]string{"boards", "create"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Lists create

func TestRun_ListsCreate_Success(t *testing.T) {
	err := runWithClient([]string{"lists", "create", "--name", "To Do", "--board-id", "b1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_ListsCreate_MissingName(t *testing.T) {
	err := runWithClient([]string{"lists", "create", "--board-id", "b1"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_ListsCreate_MissingBoardID(t *testing.T) {
	err := runWithClient([]string{"lists", "create", "--name", "List"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Cards archive

func TestRun_CardsArchive_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "archive", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsArchive_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"cards", "archive"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsUnarchive_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "unarchive", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsUnarchive_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"cards", "unarchive"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsDelete_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "delete", "--card-id", "c1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsDelete_MissingCardID(t *testing.T) {
	err := runWithClient([]string{"cards", "delete"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Cards search

func TestRun_CardsSearch_Success(t *testing.T) {
	err := runWithClient([]string{"cards", "search", "--query", "login"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsSearch_WithBoardID(t *testing.T) {
	err := runWithClient([]string{"cards", "search", "--query", "bug", "--board-id", "b1"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_CardsSearch_MissingQuery(t *testing.T) {
	err := runWithClient([]string{"cards", "search"}, testMCPClient(t))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_CardsSearch_JSON(t *testing.T) {
	err := runWithClient([]string{"cards", "search", "--query", "test", "--json"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

// Organizations

func TestRun_OrganizationsList_Success(t *testing.T) {
	err := runWithClient([]string{"orgs", "list"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRun_OrganizationsList_JSON(t *testing.T) {
	err := runWithClient([]string{"orgs", "list", "--json"}, testMCPClient(t))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}
