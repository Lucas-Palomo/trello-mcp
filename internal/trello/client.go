package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://api.trello.com/1"

type Client struct {
	APIKey     string
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey, token string, timeout time.Duration) *Client {
	return &Client{
		APIKey:     apiKey,
		Token:      token,
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) doGet(path string, v any) error {
	return c.doRequest(http.MethodGet, path, nil, v)
}

func (c *Client) doPost(path string, params url.Values, v any) error {
	return c.doRequest(http.MethodPost, path, params, v)
}

func (c *Client) doPut(path string, params url.Values, v any) error {
	return c.doRequest(http.MethodPut, path, params, v)
}

func (c *Client) doDelete(path string, params url.Values) error {
	return c.doRequest(http.MethodDelete, path, params, nil)
}

func (c *Client) doRequest(method, path string, extra url.Values, v any) error {
	req, err := http.NewRequest(method, c.BaseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Set("key", c.APIKey)
	q.Set("token", c.Token)
	for k, vals := range extra {
		for _, val := range vals {
			q.Add(k, val)
		}
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("trello request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read trello response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("trello authentication failed (401): check your API key and token")
	case http.StatusNotFound:
		return fmt.Errorf("trello resource not found (404)")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("trello returned unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if v == nil {
		return nil
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse trello response: %w", err)
	}

	return nil
}

// Navigation reads

func (c *Client) ListBoards() ([]Board, error) {
	var boards []Board
	if err := c.doGet("/members/me/boards", &boards); err != nil {
		return nil, err
	}
	return boards, nil
}

func (c *Client) ListLists(boardID string) ([]List, error) {
	var lists []List
	if err := c.doGet("/boards/"+boardID+"/lists", &lists); err != nil {
		return nil, err
	}
	return lists, nil
}

func (c *Client) ListCards(listID string) ([]Card, error) {
	var cards []Card
	if err := c.doGet("/lists/"+listID+"/cards", &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

// Detail reads

func (c *Client) GetBoard(boardID string) (Board, error) {
	var board Board
	if err := c.doGet("/boards/"+boardID, &board); err != nil {
		return Board{}, err
	}
	return board, nil
}

func (c *Client) GetCard(cardID string) (Card, error) {
	var card Card
	if err := c.doGet("/cards/"+cardID, &card); err != nil {
		return Card{}, err
	}
	return card, nil
}

// Expansion reads

func (c *Client) ListChecklists(cardID string) ([]Checklist, error) {
	var checklists []Checklist
	if err := c.doGet("/cards/"+cardID+"/checklists", &checklists); err != nil {
		return nil, err
	}
	return checklists, nil
}

func (c *Client) ListCardActions(cardID string) ([]Action, error) {
	var actions []Action
	if err := c.doGet("/cards/"+cardID+"/actions", &actions); err != nil {
		return nil, err
	}
	return actions, nil
}

func (c *Client) ListCardMembers(cardID string) ([]Member, error) {
	var members []Member
	if err := c.doGet("/cards/"+cardID+"/members", &members); err != nil {
		return nil, err
	}
	return members, nil
}

// Label mutations

func (c *Client) AddCardLabel(cardID, labelID string) error {
	var label CardLabel
	return c.doPost("/cards/"+cardID+"/idLabels", url.Values{"value": {labelID}}, &label)
}

func (c *Client) RemoveCardLabel(cardID, labelID string) error {
	return c.doDelete("/cards/"+cardID+"/idLabels/"+labelID, nil)
}

// Board and list mutations

func (c *Client) CreateBoard(name string, defaultLists bool) (Board, error) {
	params := url.Values{"name": {name}}
	if !defaultLists {
		params.Set("defaultLists", "false")
	}
	var board Board
	if err := c.doPost("/boards", params, &board); err != nil {
		return Board{}, err
	}
	return board, nil
}

func (c *Client) CreateList(name, boardID string) (List, error) {
	params := url.Values{
		"name":    {name},
		"idBoard": {boardID},
	}
	var list List
	if err := c.doPost("/lists", params, &list); err != nil {
		return List{}, err
	}
	return list, nil
}

// Organization reads

func (c *Client) ListOrganizations() ([]Organization, error) {
	var orgs []Organization
	if err := c.doGet("/members/me/organizations", &orgs); err != nil {
		return nil, err
	}
	return orgs, nil
}

// Search

type searchResponse struct {
	Cards []Card `json:"cards"`
}

func (c *Client) SearchCards(query string, boardID string) ([]Card, error) {
	params := url.Values{
		"query":      {query},
		"modelTypes": {"cards"},
	}
	if boardID != "" {
		params.Set("idBoards", boardID)
	}
	var resp searchResponse
	if err := c.doRequest(http.MethodGet, "/search", params, &resp); err != nil {
		return nil, err
	}
	return resp.Cards, nil
}

// Card archiving

func (c *Client) ArchiveCard(cardID string) (Card, error) {
	var card Card
	if err := c.doPut("/cards/"+cardID, url.Values{"closed": {"true"}}, &card); err != nil {
		return Card{}, err
	}
	return card, nil
}

func (c *Client) UnarchiveCard(cardID string) (Card, error) {
	var card Card
	if err := c.doPut("/cards/"+cardID, url.Values{"closed": {"false"}}, &card); err != nil {
		return Card{}, err
	}
	return card, nil
}

func (c *Client) DeleteCard(cardID string) error {
	return c.doDelete("/cards/"+cardID, nil)
}

// Card member mutations

func (c *Client) AddCardMember(cardID, memberID string) error {
	return c.doPost("/cards/"+cardID+"/idMembers", url.Values{"value": {memberID}}, nil)
}

func (c *Client) RemoveCardMember(cardID, memberID string) error {
	return c.doDelete("/cards/"+cardID+"/idMembers/"+memberID, nil)
}

// Label reads

func (c *Client) ListBoardLabels(boardID string) ([]Label, error) {
	var labels []Label
	if err := c.doGet("/boards/"+boardID+"/labels", &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

func (c *Client) ListCardLabels(cardID string) ([]CardLabel, error) {
	var labels []CardLabel
	if err := c.doGet("/cards/"+cardID+"/labels", &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

// Mutations

func (c *Client) CreateCard(input CreateCardInput) (Card, error) {
	params := url.Values{
		"idList": {input.ListID},
		"name":   {input.Name},
	}
	if input.Desc != "" {
		params.Set("desc", input.Desc)
	}
	if input.Due != "" {
		params.Set("due", input.Due)
	}
	var card Card
	if err := c.doPost("/cards", params, &card); err != nil {
		return Card{}, err
	}
	return card, nil
}

func (c *Client) UpdateCard(cardID string, input UpdateCardInput) (Card, error) {
	params := url.Values{}
	if input.Name != "" {
		params.Set("name", input.Name)
	}
	if input.Desc != "" {
		params.Set("desc", input.Desc)
	}
	if input.Due != "" {
		params.Set("due", input.Due)
	}
	var card Card
	if err := c.doPut("/cards/"+cardID, params, &card); err != nil {
		return Card{}, err
	}
	return card, nil
}

func (c *Client) MoveCard(cardID, listID string) (Card, error) {
	var card Card
	if err := c.doPut("/cards/"+cardID, url.Values{"idList": {listID}}, &card); err != nil {
		return Card{}, err
	}
	return card, nil
}

func (c *Client) AddComment(cardID, text string) (Action, error) {
	var action Action
	if err := c.doPost("/cards/"+cardID+"/actions/comments", url.Values{"text": {text}}, &action); err != nil {
		return Action{}, err
	}
	return action, nil
}
