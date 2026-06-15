package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"trello-mcp/internal/trello"
)

type MCPServer struct {
	Srv    *server.MCPServer
	client *trello.Client
}

func New(cfg *Config) *MCPServer {
	s := &MCPServer{
		Srv:    server.NewMCPServer("trello-mcp", "0.1.0"),
		client: cfg.TrelloClient,
	}

	// Navigation
	listBoardsTool := mcp.NewTool("trello_list_boards",
		mcp.WithDescription("Return all Trello boards accessible to the configured credentials"),
	)
	s.Srv.AddTool(listBoardsTool, s.handleListBoards)

	listListsTool := mcp.NewTool("trello_list_lists",
		mcp.WithDescription("Return the Trello lists in a board"),
		mcp.WithString("board_id", mcp.Description("Board ID to list lists from"), mcp.Required()),
	)
	s.Srv.AddTool(listListsTool, s.handleListLists)

	listCardsTool := mcp.NewTool("trello_list_cards",
		mcp.WithDescription("Return the Trello cards from a list"),
		mcp.WithString("list_id", mcp.Description("List ID to list cards from"), mcp.Required()),
	)
	s.Srv.AddTool(listCardsTool, s.handleListCards)

	// Detail reads
	getBoardTool := mcp.NewTool("trello_get_board",
		mcp.WithDescription("Get details for a Trello board"),
		mcp.WithString("board_id", mcp.Description("Board ID"), mcp.Required()),
	)
	s.Srv.AddTool(getBoardTool, s.handleGetBoard)

	getCardTool := mcp.NewTool("trello_get_card",
		mcp.WithDescription("Get details for a Trello card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(getCardTool, s.handleGetCard)

	// Expansion reads
	listChecklistsTool := mcp.NewTool("trello_list_checklists",
		mcp.WithDescription("List checklists attached to a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(listChecklistsTool, s.handleListChecklists)

	listCardActionsTool := mcp.NewTool("trello_list_card_actions",
		mcp.WithDescription("List actions (comments, updates) for a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(listCardActionsTool, s.handleListCardActions)

	listCardMembersTool := mcp.NewTool("trello_list_card_members",
		mcp.WithDescription("List members assigned to a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(listCardMembersTool, s.handleListCardMembers)

	// Mutations
	createCardTool := mcp.NewTool("trello_create_card",
		mcp.WithDescription("Create a new card in a list"),
		mcp.WithString("list_id", mcp.Description("List ID to create the card in"), mcp.Required()),
		mcp.WithString("name", mcp.Description("Card name"), mcp.Required()),
		mcp.WithString("desc", mcp.Description("Card description (optional)")),
		mcp.WithString("due", mcp.Description("Due date (RFC3339 format, optional)")),
	)
	s.Srv.AddTool(createCardTool, s.handleCreateCard)

	updateCardTool := mcp.NewTool("trello_update_card",
		mcp.WithDescription("Update editable fields on a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
		mcp.WithString("name", mcp.Description("New card name")),
		mcp.WithString("desc", mcp.Description("New card description")),
		mcp.WithString("due", mcp.Description("Due date (RFC3339 format)")),
	)
	s.Srv.AddTool(updateCardTool, s.handleUpdateCard)

	moveCardTool := mcp.NewTool("trello_move_card",
		mcp.WithDescription("Move a card to another list"),
		mcp.WithString("card_id", mcp.Description("Card ID to move"), mcp.Required()),
		mcp.WithString("list_id", mcp.Description("Target list ID"), mcp.Required()),
	)
	s.Srv.AddTool(moveCardTool, s.handleMoveCard)

	addCommentTool := mcp.NewTool("trello_add_comment",
		mcp.WithDescription("Add a comment to a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
		mcp.WithString("text", mcp.Description("Comment text"), mcp.Required()),
	)
	s.Srv.AddTool(addCommentTool, s.handleAddComment)

	// Label reads
	listBoardLabelsTool := mcp.NewTool("trello_list_board_labels",
		mcp.WithDescription("List all labels defined on a board"),
		mcp.WithString("board_id", mcp.Description("Board ID"), mcp.Required()),
	)
	s.Srv.AddTool(listBoardLabelsTool, s.handleListBoardLabels)

	listCardLabelsTool := mcp.NewTool("trello_list_card_labels",
		mcp.WithDescription("List labels applied to a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(listCardLabelsTool, s.handleListCardLabels)

	// Label mutations
	addCardLabelTool := mcp.NewTool("trello_add_card_label",
		mcp.WithDescription("Attach a label to a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
		mcp.WithString("label_id", mcp.Description("Label ID"), mcp.Required()),
	)
	s.Srv.AddTool(addCardLabelTool, s.handleAddCardLabel)

	removeCardLabelTool := mcp.NewTool("trello_remove_card_label",
		mcp.WithDescription("Remove a label from a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
		mcp.WithString("label_id", mcp.Description("Label ID"), mcp.Required()),
	)
	s.Srv.AddTool(removeCardLabelTool, s.handleRemoveCardLabel)

	// Card member mutations
	addCardMemberTool := mcp.NewTool("trello_add_card_member",
		mcp.WithDescription("Assign a member to a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
		mcp.WithString("member_id", mcp.Description("Member ID"), mcp.Required()),
	)
	s.Srv.AddTool(addCardMemberTool, s.handleAddCardMember)

	removeCardMemberTool := mcp.NewTool("trello_remove_card_member",
		mcp.WithDescription("Remove a member from a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
		mcp.WithString("member_id", mcp.Description("Member ID"), mcp.Required()),
	)
	s.Srv.AddTool(removeCardMemberTool, s.handleRemoveCardMember)

	// Board and list mutations
	createBoardTool := mcp.NewTool("trello_create_board",
		mcp.WithDescription("Create a new Trello board"),
		mcp.WithString("name", mcp.Description("Board name"), mcp.Required()),
		mcp.WithBoolean("default_lists", mcp.Description("Include default lists (optional, default true)")),
	)
	s.Srv.AddTool(createBoardTool, s.handleCreateBoard)

	createListTool := mcp.NewTool("trello_create_list",
		mcp.WithDescription("Create a new list in a board"),
		mcp.WithString("name", mcp.Description("List name"), mcp.Required()),
		mcp.WithString("board_id", mcp.Description("Board ID"), mcp.Required()),
	)
	s.Srv.AddTool(createListTool, s.handleCreateList)

	// Card archiving
	archiveCardTool := mcp.NewTool("trello_archive_card",
		mcp.WithDescription("Archive (close) a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(archiveCardTool, s.handleArchiveCard)

	unarchiveCardTool := mcp.NewTool("trello_unarchive_card",
		mcp.WithDescription("Unarchive (re-open) a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(unarchiveCardTool, s.handleUnarchiveCard)

	deleteCardTool := mcp.NewTool("trello_delete_card",
		mcp.WithDescription("Permanently delete a card"),
		mcp.WithString("card_id", mcp.Description("Card ID"), mcp.Required()),
	)
	s.Srv.AddTool(deleteCardTool, s.handleDeleteCard)

	// Organization reads
	listOrgsTool := mcp.NewTool("trello_list_organizations",
		mcp.WithDescription("List Trello organizations (workspaces) for the configured account"),
	)
	s.Srv.AddTool(listOrgsTool, s.handleListOrganizations)

	// Search
	searchCardsTool := mcp.NewTool("trello_search_cards",
		mcp.WithDescription("Search for cards by text query"),
		mcp.WithString("query", mcp.Description("Search query"), mcp.Required()),
		mcp.WithString("board_id", mcp.Description("Board ID to scope search (optional)")),
	)
	s.Srv.AddTool(searchCardsTool, s.handleSearchCards)

	return s
}

type Config struct {
	TrelloClient *trello.Client
}

func (s *MCPServer) Serve() error {
	return server.ServeStdio(s.Srv)
}

// Navigation handlers

func (s *MCPServer) handleListBoards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boards, err := s.client.ListBoards()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list boards: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.BoardListResult{Boards: boards})
}

func (s *MCPServer) handleListLists(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := mcp.ParseString(req, "board_id", "")
	if boardID == "" {
		return mcp.NewToolResultError("board_id is required"), nil
	}
	lists, err := s.client.ListLists(boardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list lists: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.ListListResult{Lists: lists})
}

func (s *MCPServer) handleListCards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	listID := mcp.ParseString(req, "list_id", "")
	if listID == "" {
		return mcp.NewToolResultError("list_id is required"), nil
	}
	cards, err := s.client.ListCards(listID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list cards: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.CardListResult{Cards: cards})
}

// Detail handlers

func (s *MCPServer) handleGetBoard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := mcp.ParseString(req, "board_id", "")
	if boardID == "" {
		return mcp.NewToolResultError("board_id is required"), nil
	}
	board, err := s.client.GetBoard(boardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get board: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.GetBoardResult{Board: board})
}

func (s *MCPServer) handleGetCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	card, err := s.client.GetCard(cardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get card: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.GetCardResult{Card: card})
}

// Expansion handlers

func (s *MCPServer) handleListChecklists(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	checklists, err := s.client.ListChecklists(cardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list checklists: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.ChecklistListResult{Checklists: checklists})
}

func (s *MCPServer) handleListCardActions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	actions, err := s.client.ListCardActions(cardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list card actions: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.ActionListResult{Actions: actions})
}

func (s *MCPServer) handleListCardMembers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	members, err := s.client.ListCardMembers(cardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list card members: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.MemberListResult{Members: members})
}

// Mutation handlers

func (s *MCPServer) handleCreateCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	listID := mcp.ParseString(req, "list_id", "")
	name := mcp.ParseString(req, "name", "")
	desc := mcp.ParseString(req, "desc", "")

	if listID == "" {
		return mcp.NewToolResultError("list_id is required"), nil
	}
	if name == "" {
		return mcp.NewToolResultError("name is required"), nil
	}

	due := mcp.ParseString(req, "due", "")

	card, err := s.client.CreateCard(trello.CreateCardInput{
		ListID: listID,
		Name:   name,
		Desc:   desc,
		Due:    due,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create card: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.CreateCardResult{Card: card})
}

func (s *MCPServer) handleUpdateCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}

	input := trello.UpdateCardInput{
		Name: mcp.ParseString(req, "name", ""),
		Desc: mcp.ParseString(req, "desc", ""),
		Due:  mcp.ParseString(req, "due", ""),
	}

	if input.Name == "" && input.Desc == "" && input.Due == "" {
		return mcp.NewToolResultError("at least one of name, desc, or due must be provided"), nil
	}

	card, err := s.client.UpdateCard(cardID, input)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update card: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.UpdateCardResult{Card: card})
}

func (s *MCPServer) handleMoveCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	listID := mcp.ParseString(req, "list_id", "")

	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	if listID == "" {
		return mcp.NewToolResultError("list_id is required"), nil
	}

	card, err := s.client.MoveCard(cardID, listID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to move card: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.MoveCardResult{Card: card})
}

func (s *MCPServer) handleAddComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	text := mcp.ParseString(req, "text", "")

	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	if text == "" {
		return mcp.NewToolResultError("text is required"), nil
	}

	action, err := s.client.AddComment(cardID, text)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to add comment: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.AddCommentResult{Action: action})
}

// Label handlers

func (s *MCPServer) handleListBoardLabels(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := mcp.ParseString(req, "board_id", "")
	if boardID == "" {
		return mcp.NewToolResultError("board_id is required"), nil
	}
	labels, err := s.client.ListBoardLabels(boardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list board labels: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.BoardLabelListResult{Labels: labels})
}

func (s *MCPServer) handleListCardLabels(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	labels, err := s.client.ListCardLabels(cardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list card labels: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.CardLabelListResult{Labels: labels})
}

// Label mutation handlers

func (s *MCPServer) handleAddCardLabel(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	labelID := mcp.ParseString(req, "label_id", "")

	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	if labelID == "" {
		return mcp.NewToolResultError("label_id is required"), nil
	}

	if err := s.client.AddCardLabel(cardID, labelID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to add label: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.LabelActionResult{Success: true})
}

func (s *MCPServer) handleRemoveCardLabel(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	labelID := mcp.ParseString(req, "label_id", "")

	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	if labelID == "" {
		return mcp.NewToolResultError("label_id is required"), nil
	}

	if err := s.client.RemoveCardLabel(cardID, labelID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove label: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.LabelActionResult{Success: true})
}

// Card member mutation handlers

func (s *MCPServer) handleAddCardMember(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	memberID := mcp.ParseString(req, "member_id", "")

	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	if memberID == "" {
		return mcp.NewToolResultError("member_id is required"), nil
	}

	if err := s.client.AddCardMember(cardID, memberID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to add member: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.LabelActionResult{Success: true})
}

func (s *MCPServer) handleRemoveCardMember(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	memberID := mcp.ParseString(req, "member_id", "")

	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}
	if memberID == "" {
		return mcp.NewToolResultError("member_id is required"), nil
	}

	if err := s.client.RemoveCardMember(cardID, memberID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove member: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.LabelActionResult{Success: true})
}

// Board and list mutation handlers

func (s *MCPServer) handleCreateBoard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := mcp.ParseString(req, "name", "")
	if name == "" {
		return mcp.NewToolResultError("name is required"), nil
	}
	defaultLists := mcp.ParseBoolean(req, "default_lists", true)

	board, err := s.client.CreateBoard(name, defaultLists)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create board: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.CreateBoardResult{Board: board})
}

func (s *MCPServer) handleCreateList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := mcp.ParseString(req, "name", "")
	boardID := mcp.ParseString(req, "board_id", "")

	if name == "" {
		return mcp.NewToolResultError("name is required"), nil
	}
	if boardID == "" {
		return mcp.NewToolResultError("board_id is required"), nil
	}

	list, err := s.client.CreateList(name, boardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create list: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.CreateListResult{List: list})
}

// Card archiving handlers

func (s *MCPServer) handleArchiveCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}

	card, err := s.client.ArchiveCard(cardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to archive card: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.ArchiveCardResult{Card: card})
}

func (s *MCPServer) handleUnarchiveCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}

	card, err := s.client.UnarchiveCard(cardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to unarchive card: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.UnarchiveCardResult{Card: card})
}

func (s *MCPServer) handleDeleteCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cardID := mcp.ParseString(req, "card_id", "")
	if cardID == "" {
		return mcp.NewToolResultError("card_id is required"), nil
	}

	if err := s.client.DeleteCard(cardID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete card: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.LabelActionResult{Success: true})
}

// Organization handler

func (s *MCPServer) handleListOrganizations(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgs, err := s.client.ListOrganizations()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list organizations: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.OrganizationListResult{Organizations: orgs})
}

// Search handler

func (s *MCPServer) handleSearchCards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := mcp.ParseString(req, "query", "")
	if query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}
	boardID := mcp.ParseString(req, "board_id", "")

	cards, err := s.client.SearchCards(query, boardID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search cards: %v", err)), nil
	}
	return mcp.NewToolResultJSON(trello.SearchCardsResult{Cards: cards})
}
