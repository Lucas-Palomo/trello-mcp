package trello

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testClient(handler http.HandlerFunc) *Client {
	srv := httptest.NewServer(handler)
	return &Client{
		APIKey:     "test-key",
		Token:      "test-token",
		BaseURL:    srv.URL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func testOK(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

// --- ListBoards ---

func TestListBoards_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "b1", "name": "Board One", "url": "https://trello.com/b/b1"},
		})
	})
	boards, err := c.ListBoards()
	if err != nil {
		t.Fatal(err)
	}
	if len(boards) != 1 || boards[0].ID != "b1" {
		t.Fatalf("unexpected: %+v", boards)
	}
}

func TestListBoards_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.ListBoards()
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListCards ---

func TestListCards_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "c1", "name": "Card One", "desc": "Desc 1"},
		})
	})
	cards, err := c.ListCards("list123")
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 || cards[0].ID != "c1" {
		t.Fatalf("unexpected: %+v", cards)
	}
}

func TestListCards_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.ListCards("list123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListCards_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ListCards("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListCards_ServerError(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	_, err := c.ListCards("list123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListCards_InvalidJSON(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not json`))
	})
	_, err := c.ListCards("list123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListCards_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))
	defer srv.Close()
	c := &Client{
		APIKey: "key", Token: "token", BaseURL: srv.URL,
		HTTPClient: &http.Client{Timeout: 1 * time.Millisecond},
	}
	_, err := c.ListCards("list123")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListLists ---

func TestListLists_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "l1", "name": "List One"},
		})
	})
	lists, err := c.ListLists("board123")
	if err != nil {
		t.Fatal(err)
	}
	if len(lists) != 1 || lists[0].ID != "l1" {
		t.Fatalf("unexpected: %+v", lists)
	}
}

func TestListLists_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.ListLists("board123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListLists_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ListLists("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- GetBoard ---

func TestGetBoard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, map[string]any{"id": "b1", "name": "Board 1", "url": "https://trello.com/b/b1"})
	})
	board, err := c.GetBoard("b1")
	if err != nil {
		t.Fatal(err)
	}
	if board.ID != "b1" || board.Name != "Board 1" {
		t.Fatalf("unexpected: %+v", board)
	}
}

func TestGetBoard_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.GetBoard("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- GetCard ---

func TestGetCard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, map[string]any{"id": "c1", "name": "Card 1", "desc": "A card"})
	})
	card, err := c.GetCard("c1")
	if err != nil {
		t.Fatal(err)
	}
	if card.ID != "c1" || card.Name != "Card 1" {
		t.Fatalf("unexpected: %+v", card)
	}
}

func TestGetCard_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.GetCard("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListChecklists ---

func TestListChecklists_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "chk1", "name": "Checklist 1"},
		})
	})
	cl, err := c.ListChecklists("card123")
	if err != nil {
		t.Fatal(err)
	}
	if len(cl) != 1 || cl[0].ID != "chk1" {
		t.Fatalf("unexpected: %+v", cl)
	}
}

func TestListChecklists_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ListChecklists("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListCardActions ---

func TestListCardActions_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "a1", "type": "commentCard", "text": "nice"},
		})
	})
	actions, err := c.ListCardActions("card123")
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 1 || actions[0].ID != "a1" {
		t.Fatalf("unexpected: %+v", actions)
	}
}

func TestListCardActions_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ListCardActions("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListCardMembers ---

func TestListCardMembers_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "m1", "full_name": "Jane Doe", "username": "jane"},
		})
	})
	members, err := c.ListCardMembers("card123")
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 1 || members[0].ID != "m1" {
		t.Fatalf("unexpected: %+v", members)
	}
}

func TestListCardMembers_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ListCardMembers("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- CreateCard ---

func TestCreateCard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("idList") != "l1" || r.URL.Query().Get("name") != "New Card" {
			t.Errorf("unexpected params: %v", r.URL.Query())
		}
		testOK(w, map[string]any{"id": "c1", "name": "New Card", "url": "https://trello.com/c/c1"})
	})
	card, err := c.CreateCard(CreateCardInput{ListID: "l1", Name: "New Card"})
	if err != nil {
		t.Fatal(err)
	}
	if card.ID != "c1" || card.Name != "New Card" {
		t.Fatalf("unexpected: %+v", card)
	}
}

func TestCreateCard_WithDesc(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("desc") != "A description" {
			t.Errorf("unexpected desc: %s", r.URL.Query().Get("desc"))
		}
		testOK(w, map[string]any{"id": "c1", "name": "Card"})
	})
	_, err := c.CreateCard(CreateCardInput{ListID: "l1", Name: "Card", Desc: "A description"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateCard_WithDue(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("due") != "2026-07-01T12:00:00Z" {
			t.Errorf("unexpected due: %s", r.URL.Query().Get("due"))
		}
		testOK(w, map[string]any{"id": "c1", "name": "Card", "due": "2026-07-01T12:00:00Z"})
	})
	card, err := c.CreateCard(CreateCardInput{ListID: "l1", Name: "Card", Due: "2026-07-01T12:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	if card.Due != "2026-07-01T12:00:00Z" {
		t.Fatalf("unexpected due: %s", card.Due)
	}
}

func TestCreateCard_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.CreateCard(CreateCardInput{ListID: "l1", Name: "Card"})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- UpdateCard ---

func TestUpdateCard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Query().Get("name") != "Updated" {
			t.Errorf("unexpected name: %s", r.URL.Query().Get("name"))
		}
		testOK(w, map[string]any{"id": "c1", "name": "Updated"})
	})
	card, err := c.UpdateCard("c1", UpdateCardInput{Name: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if card.ID != "c1" || card.Name != "Updated" {
		t.Fatalf("unexpected: %+v", card)
	}
}

func TestUpdateCard_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.UpdateCard("c1", UpdateCardInput{Name: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- MoveCard ---

func TestMoveCard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Query().Get("idList") != "l2" {
			t.Errorf("unexpected idList: %s", r.URL.Query().Get("idList"))
		}
		testOK(w, map[string]any{"id": "c1", "name": "Moved Card"})
	})
	card, err := c.MoveCard("c1", "l2")
	if err != nil {
		t.Fatal(err)
	}
	if card.ID != "c1" {
		t.Fatalf("unexpected: %+v", card)
	}
}

func TestMoveCard_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.MoveCard("invalid", "l2")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- AddComment ---

func TestAddComment_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("text") != "Great card!" {
			t.Errorf("unexpected text: %s", r.URL.Query().Get("text"))
		}
		testOK(w, map[string]any{"id": "a1", "type": "commentCard", "text": "Great card!"})
	})
	action, err := c.AddComment("c1", "Great card!")
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != "a1" || action.Text != "Great card!" {
		t.Fatalf("unexpected: %+v", action)
	}
}

func TestAddComment_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.AddComment("c1", "text")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListBoardLabels ---

func TestListBoardLabels_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "lb1", "idBoard": "b1", "name": "urgent", "color": "red"},
			{"id": "lb2", "idBoard": "b1", "name": "bug", "color": "orange"},
		})
	})
	labels, err := c.ListBoardLabels("b1")
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 2 || labels[0].ID != "lb1" {
		t.Fatalf("unexpected: %+v", labels)
	}
}

func TestListBoardLabels_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ListBoardLabels("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListCardLabels ---

func TestListCardLabels_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "lb1", "idBoard": "b1", "name": "urgent", "color": "red"},
		})
	})
	labels, err := c.ListCardLabels("c1")
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 || labels[0].Name != "urgent" {
		t.Fatalf("unexpected: %+v", labels)
	}
}

func TestListCardLabels_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ListCardLabels("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- AddCardLabel ---

func TestAddCardLabel_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("value") != "lb1" {
			t.Errorf("unexpected value: %s", r.URL.Query().Get("value"))
		}
		testOK(w, map[string]any{"id": "lb1", "name": "urgent", "color": "red"})
	})
	err := c.AddCardLabel("c1", "lb1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddCardLabel_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	err := c.AddCardLabel("invalid", "lb1")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- RemoveCardLabel ---

func TestRemoveCardLabel_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	err := c.RemoveCardLabel("c1", "lb1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveCardLabel_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	err := c.RemoveCardLabel("invalid", "lb1")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- AddCardMember ---

func TestAddCardMember_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("value") != "m1" {
			t.Errorf("unexpected value: %s", r.URL.Query().Get("value"))
		}
		w.WriteHeader(http.StatusOK)
	})
	err := c.AddCardMember("c1", "m1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddCardMember_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	err := c.AddCardMember("invalid", "m1")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- RemoveCardMember ---

func TestRemoveCardMember_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	err := c.RemoveCardMember("c1", "m1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveCardMember_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	err := c.RemoveCardMember("invalid", "m1")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- CreateBoard ---

func TestCreateBoard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("name") != "My Board" {
			t.Errorf("unexpected name: %s", r.URL.Query().Get("name"))
		}
		testOK(w, map[string]any{"id": "b1", "name": "My Board", "url": "https://trello.com/b/b1"})
	})
	board, err := c.CreateBoard("My Board", true)
	if err != nil {
		t.Fatal(err)
	}
	if board.ID != "b1" || board.Name != "My Board" {
		t.Fatalf("unexpected: %+v", board)
	}
}

func TestCreateBoard_NoDefaultLists(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("defaultLists") != "false" {
			t.Errorf("expected defaultLists=false, got %s", r.URL.Query().Get("defaultLists"))
		}
		testOK(w, map[string]any{"id": "b1", "name": "Board"})
	})
	_, err := c.CreateBoard("Board", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateBoard_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.CreateBoard("Board", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- CreateList ---

func TestCreateList_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("name") != "To Do" || r.URL.Query().Get("idBoard") != "b1" {
			t.Errorf("unexpected params: %v", r.URL.Query())
		}
		testOK(w, map[string]any{"id": "l1", "name": "To Do"})
	})
	list, err := c.CreateList("To Do", "b1")
	if err != nil {
		t.Fatal(err)
	}
	if list.ID != "l1" || list.Name != "To Do" {
		t.Fatalf("unexpected: %+v", list)
	}
}

func TestCreateList_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.CreateList("List", "invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ArchiveCard ---

func TestArchiveCard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Query().Get("closed") != "true" {
			t.Errorf("expected closed=true, got %s", r.URL.Query().Get("closed"))
		}
		testOK(w, map[string]any{"id": "c1", "closed": true})
	})
	card, err := c.ArchiveCard("c1")
	if err != nil {
		t.Fatal(err)
	}
	if !card.Closed {
		t.Fatal("expected closed=true")
	}
}

func TestArchiveCard_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.ArchiveCard("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- UnarchiveCard ---

func TestUnarchiveCard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Query().Get("closed") != "false" {
			t.Errorf("expected closed=false, got %s", r.URL.Query().Get("closed"))
		}
		testOK(w, map[string]any{"id": "c1", "closed": false})
	})
	card, err := c.UnarchiveCard("c1")
	if err != nil {
		t.Fatal(err)
	}
	if card.Closed {
		t.Fatal("expected closed=false")
	}
}

// --- DeleteCard ---

func TestDeleteCard_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	err := c.DeleteCard("c1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteCard_NotFound(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	err := c.DeleteCard("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- SearchCards ---

func TestSearchCards_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("query") != "login" {
			t.Errorf("unexpected query: %s", r.URL.Query().Get("query"))
		}
		if r.URL.Query().Get("modelTypes") != "cards" {
			t.Errorf("unexpected modelTypes: %s", r.URL.Query().Get("modelTypes"))
		}
		testOK(w, map[string]any{
			"cards": []map[string]any{
				{"id": "c1", "name": "Fix login", "idList": "l1", "idBoard": "b1"},
			},
		})
	})
	cards, err := c.SearchCards("login", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 || cards[0].Name != "Fix login" {
		t.Fatalf("unexpected: %+v", cards)
	}
	if cards[0].IDBoard != "b1" {
		t.Fatalf("expected idBoard=b1, got %q", cards[0].IDBoard)
	}
}

func TestSearchCards_WithBoardID(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("idBoards") != "b1" {
			t.Errorf("unexpected idBoards: %s", r.URL.Query().Get("idBoards"))
		}
		testOK(w, map[string]any{"cards": []map[string]any{}})
	})
	_, err := c.SearchCards("bug", "b1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchCards_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.SearchCards("test", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListOrganizations ---

func TestListOrganizations_Success(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		testOK(w, []map[string]any{
			{"id": "o1", "name": "acmecorp", "displayName": "Acme Corp", "url": "https://trello.com/acmecorp"},
		})
	})
	orgs, err := c.ListOrganizations()
	if err != nil {
		t.Fatal(err)
	}
	if len(orgs) != 1 || orgs[0].Name != "acmecorp" {
		t.Fatalf("unexpected: %+v", orgs)
	}
}

func TestListOrganizations_Unauthorized(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := c.ListOrganizations()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListBoards_AuthQueryParams(t *testing.T) {
	c := testClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-key" || r.URL.Query().Get("token") != "test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		testOK(w, []map[string]any{})
	})
	_, err := c.ListBoards()
	if err != nil {
		t.Fatal(err)
	}
}
