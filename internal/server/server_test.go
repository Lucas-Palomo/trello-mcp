package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"trello-mcp/internal/trello"
)

func testClient(handler http.HandlerFunc) *trello.Client {
	srv := httptest.NewServer(handler)
	return &trello.Client{
		APIKey:     "key",
		Token:      "token",
		BaseURL:    srv.URL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func args(m map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{Arguments: m},
	}
}

// --- ListBoards ---

func TestHandleListBoards_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"b1","name":"Board 1","url":"https://trello.com/b/b1"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListBoards(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.BoardListResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if len(resp.Boards) != 1 {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestHandleListBoards_Error(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListBoards(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleListBoards_InvalidJSON(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListBoards(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListLists ---

func TestHandleListLists_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"l1","name":"List 1"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListLists(context.Background(), args(map[string]any{"board_id": "b1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleListLists_MissingBoardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleListLists(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleListLists_NotFound(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListLists(context.Background(), args(map[string]any{"board_id": "invalid"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListCards ---

func TestHandleListCards_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"c1","name":"Card 1","desc":"A card"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListCards(context.Background(), args(map[string]any{"list_id": "l1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleListCards_MissingListID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleListCards(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleListCards_NotFound(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListCards(context.Background(), args(map[string]any{"list_id": "invalid"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- GetBoard ---

func TestHandleGetBoard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"b1","name":"Board 1","url":"https://trello.com/b/b1"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleGetBoard(context.Background(), args(map[string]any{"board_id": "b1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.GetBoardResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if resp.Board.ID != "b1" {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestHandleGetBoard_MissingBoardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleGetBoard(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleGetBoard_NotFound(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleGetBoard(context.Background(), args(map[string]any{"board_id": "invalid"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- GetCard ---

func TestHandleGetCard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"c1","name":"Card 1","desc":"A card","due":"2026-07-01T12:00:00Z","idList":"l1","labels":[{"id":"lb1","name":"urgent","color":"red"}]}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleGetCard(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.GetCardResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if resp.Card.ID != "c1" {
		t.Fatalf("unexpected: %+v", resp)
	}
	if resp.Card.Due != "2026-07-01T12:00:00Z" {
		t.Fatalf("expected due date, got: %s", resp.Card.Due)
	}
	if len(resp.Card.Labels) != 1 || resp.Card.Labels[0].Name != "urgent" {
		t.Fatalf("unexpected labels: %+v", resp.Card.Labels)
	}
}

func TestHandleGetCard_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleGetCard(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleGetCard_NotFound(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleGetCard(context.Background(), args(map[string]any{"card_id": "invalid"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListChecklists ---

func TestHandleListChecklists_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"chk1","name":"Checklist 1"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListChecklists(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleListChecklists_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleListChecklists(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListCardActions ---

func TestHandleListCardActions_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"a1","type":"commentCard","text":"nice"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListCardActions(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleListCardActions_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleListCardActions(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListCardMembers ---

func TestHandleListCardMembers_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"m1","full_name":"Jane","username":"jane"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListCardMembers(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleListCardMembers_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleListCardMembers(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- CreateCard ---

func TestHandleCreateCard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"c1","name":"New Card","url":"https://trello.com/c/c1"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleCreateCard(context.Background(), args(map[string]any{
		"list_id": "l1", "name": "New Card",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleCreateCard_WithDue(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("due") != "2026-07-01T12:00:00Z" {
			t.Errorf("unexpected due: %s", r.URL.Query().Get("due"))
		}
		w.Write([]byte(`{"id":"c1","name":"New Card","due":"2026-07-01T12:00:00Z"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleCreateCard(context.Background(), args(map[string]any{
		"list_id": "l1", "name": "New Card", "due": "2026-07-01T12:00:00Z",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.CreateCardResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if resp.Card.Due != "2026-07-01T12:00:00Z" {
		t.Fatalf("expected due, got: %s", resp.Card.Due)
	}
}

func TestHandleCreateCard_MissingListID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleCreateCard(context.Background(), args(map[string]any{"name": "Card"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleCreateCard_MissingName(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleCreateCard(context.Background(), args(map[string]any{"list_id": "l1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- UpdateCard ---

func TestHandleUpdateCard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"c1","name":"Updated","url":"https://trello.com/c/c1"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleUpdateCard(context.Background(), args(map[string]any{
		"card_id": "c1", "name": "Updated",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleUpdateCard_ExtendedFields(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"c1","name":"Updated","due":"2026-08-01T12:00:00Z","labels":[{"id":"lb1","name":"urgent","color":"red"}],"idMembers":["m1"],"idList":"l2"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleUpdateCard(context.Background(), args(map[string]any{
		"card_id": "c1", "due": "2026-08-01T12:00:00Z",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.UpdateCardResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if resp.Card.Due != "2026-08-01T12:00:00Z" {
		t.Fatalf("expected due, got: %s", resp.Card.Due)
	}
	if len(resp.Card.Labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(resp.Card.Labels))
	}
	if len(resp.Card.IDMembers) != 1 {
		t.Fatalf("expected 1 member, got %d", len(resp.Card.IDMembers))
	}
	if resp.Card.IDList != "l2" {
		t.Fatalf("expected idList l2, got %s", resp.Card.IDList)
	}
}

func TestHandleUpdateCard_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleUpdateCard(context.Background(), args(map[string]any{"name": "x"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleUpdateCard_NoFields(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleUpdateCard(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- MoveCard ---

func TestHandleMoveCard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"c1","name":"Moved","url":"https://trello.com/c/c1"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleMoveCard(context.Background(), args(map[string]any{
		"card_id": "c1", "list_id": "l2",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleMoveCard_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleMoveCard(context.Background(), args(map[string]any{"list_id": "l2"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleMoveCard_MissingListID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleMoveCard(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- AddComment ---

func TestHandleAddComment_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"a1","type":"commentCard","text":"Nice!"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleAddComment(context.Background(), args(map[string]any{
		"card_id": "c1", "text": "Nice!",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleAddComment_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleAddComment(context.Background(), args(map[string]any{"text": "hi"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleAddComment_MissingText(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleAddComment(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListBoardLabels ---

func TestHandleListBoardLabels_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"lb1","idBoard":"b1","name":"urgent","color":"red"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListBoardLabels(context.Background(), args(map[string]any{"board_id": "b1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.BoardLabelListResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if len(resp.Labels) != 1 || resp.Labels[0].Name != "urgent" {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestHandleListBoardLabels_MissingBoardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleListBoardLabels(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListCardLabels ---

func TestHandleListCardLabels_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"lb1","idBoard":"b1","name":"urgent","color":"red"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListCardLabels(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.CardLabelListResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if len(resp.Labels) != 1 || resp.Labels[0].Name != "urgent" {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestHandleListCardLabels_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleListCardLabels(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- AddCardLabel ---

func TestHandleAddCardLabel_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"lb1","name":"urgent","color":"red"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleAddCardLabel(context.Background(), args(map[string]any{"card_id": "c1", "label_id": "lb1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleAddCardLabel_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleAddCardLabel(context.Background(), args(map[string]any{"label_id": "lb1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleAddCardLabel_MissingLabelID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleAddCardLabel(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- RemoveCardLabel ---

func TestHandleRemoveCardLabel_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleRemoveCardLabel(context.Background(), args(map[string]any{"card_id": "c1", "label_id": "lb1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleRemoveCardLabel_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleRemoveCardLabel(context.Background(), args(map[string]any{"label_id": "lb1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleRemoveCardLabel_MissingLabelID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleRemoveCardLabel(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- AddCardMember ---

func TestHandleAddCardMember_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleAddCardMember(context.Background(), args(map[string]any{"card_id": "c1", "member_id": "m1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleAddCardMember_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleAddCardMember(context.Background(), args(map[string]any{"member_id": "m1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleAddCardMember_MissingMemberID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleAddCardMember(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- RemoveCardMember ---

func TestHandleRemoveCardMember_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleRemoveCardMember(context.Background(), args(map[string]any{"card_id": "c1", "member_id": "m1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleRemoveCardMember_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleRemoveCardMember(context.Background(), args(map[string]any{"member_id": "m1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleRemoveCardMember_MissingMemberID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleRemoveCardMember(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- CreateBoard ---

func TestHandleCreateBoard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"b1","name":"My Board","url":"https://trello.com/b/b1"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleCreateBoard(context.Background(), args(map[string]any{"name": "My Board"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.CreateBoardResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if resp.Board.ID != "b1" || resp.Board.Name != "My Board" {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestHandleCreateBoard_MissingName(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleCreateBoard(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- CreateList ---

func TestHandleCreateList_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"l1","name":"To Do"}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleCreateList(context.Background(), args(map[string]any{"name": "To Do", "board_id": "b1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.CreateListResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if resp.List.ID != "l1" || resp.List.Name != "To Do" {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestHandleCreateList_MissingName(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleCreateList(context.Background(), args(map[string]any{"board_id": "b1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

func TestHandleCreateList_MissingBoardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleCreateList(context.Background(), args(map[string]any{"name": "List"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ArchiveCard ---

func TestHandleArchiveCard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"c1","closed":true}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleArchiveCard(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.ArchiveCardResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if !resp.Card.Closed {
		t.Fatal("expected closed=true")
	}
}

func TestHandleArchiveCard_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleArchiveCard(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- UnarchiveCard ---

func TestHandleUnarchiveCard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"c1","closed":false}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleUnarchiveCard(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.UnarchiveCardResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if resp.Card.Closed {
		t.Fatal("expected closed=false")
	}
}

func TestHandleUnarchiveCard_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleUnarchiveCard(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- DeleteCard ---

func TestHandleDeleteCard_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleDeleteCard(context.Background(), args(map[string]any{"card_id": "c1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleDeleteCard_MissingCardID(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleDeleteCard(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- SearchCards ---

func TestHandleSearchCards_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"cards":[{"id":"c1","name":"Fix login","idList":"l1","idBoard":"b1"}]}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleSearchCards(context.Background(), args(map[string]any{"query": "login"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.SearchCardsResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if len(resp.Cards) != 1 || resp.Cards[0].Name != "Fix login" {
		t.Fatalf("unexpected: %+v", resp)
	}
	if resp.Cards[0].IDBoard != "b1" {
		t.Fatalf("expected idBoard=b1, got %q", resp.Cards[0].IDBoard)
	}
}

func TestHandleSearchCards_WithBoardID(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"cards":[]}`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleSearchCards(context.Background(), args(map[string]any{"query": "bug", "board_id": "b1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
}

func TestHandleSearchCards_MissingQuery(t *testing.T) {
	s := New(&Config{TrelloClient: testClient(nil)})
	result, err := s.handleSearchCards(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}

// --- ListOrganizations ---

func TestHandleListOrganizations_Success(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"o1","name":"acmecorp","displayName":"Acme Corp","url":"https://trello.com/acmecorp"}]`))
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListOrganizations(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content[0].(mcp.TextContent).Text)
	}
	var resp trello.OrganizationListResult
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &resp)
	if len(resp.Organizations) != 1 || resp.Organizations[0].Name != "acmecorp" {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestHandleListOrganizations_Error(t *testing.T) {
	client := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	s := New(&Config{TrelloClient: client})
	result, err := s.handleListOrganizations(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error")
	}
}
