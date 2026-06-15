package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func printJSON(result any) error {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func printToolsList(result *mcp.ListToolsResult) error {
	for _, t := range result.Tools {
		fmt.Printf("  %s\n", t.Name)
		if t.Description != "" {
			fmt.Printf("    %s\n", t.Description)
		}
	}
	return nil
}

func printBoardsList(result *mcp.CallToolResult) error {
	var data struct {
		Boards []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"boards"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, b := range data.Boards {
		fmt.Printf("  %-40s %s\n", b.Name, b.ID)
	}
	return nil
}

func printBoardDetail(result *mcp.CallToolResult) error {
	var data struct {
		Board struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"board"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	fmt.Printf("  %s\n", data.Board.Name)
	fmt.Printf("  ID: %s\n", data.Board.ID)
	fmt.Printf("  URL: %s\n", data.Board.URL)
	return nil
}

func printListsList(result *mcp.CallToolResult) error {
	var data struct {
		Lists []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"lists"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, l := range data.Lists {
		fmt.Printf("  %-40s %s\n", l.Name, l.ID)
	}
	return nil
}

type cardPrint struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	URL         string           `json:"url"`
	Desc        string           `json:"desc"`
	Due         string           `json:"due"`
	Closed      bool             `json:"closed"`
	IDList      string           `json:"idList"`
	IDMembers   []string         `json:"idMembers"`
	Labels      []cardLabelPrint `json:"labels"`
}

type cardLabelPrint struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func printCardsList(result *mcp.CallToolResult) error {
	var data struct {
		Cards []cardPrint `json:"cards"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, c := range data.Cards {
		desc := strings.ReplaceAll(c.Desc, "\n", " ")
		if len(desc) > 60 {
			desc = desc[:60] + "..."
		}
		fmt.Printf("  %-40s %s\n", c.Name, c.ID)
		if desc != "" {
			fmt.Printf("    %s\n", desc)
		}
		if c.Due != "" {
			fmt.Printf("    Due: %s\n", c.Due)
		}
		if len(c.Labels) > 0 {
			parts := make([]string, 0, len(c.Labels))
			for _, l := range c.Labels {
				parts = append(parts, l.Name)
			}
			fmt.Printf("    Labels: %s\n", strings.Join(parts, ", "))
		}
	}
	return nil
}

func printCardDetail(result *mcp.CallToolResult) error {
	text := resultText(result)

	var card struct {
		Card *cardPrint `json:"card"`
	}
	if json.Unmarshal([]byte(text), &card) == nil && card.Card != nil {
		c := card.Card
		fmt.Printf("  %s\n", c.Name)
		fmt.Printf("  ID: %s\n", c.ID)
		fmt.Printf("  URL: %s\n", c.URL)
		if c.Desc != "" {
			fmt.Printf("  Desc: %s\n", c.Desc)
		}
		if c.Due != "" {
			fmt.Printf("  Due: %s\n", c.Due)
		}
		if c.Closed {
			fmt.Println("  Status: archived")
		}
		if len(c.Labels) > 0 {
			parts := make([]string, 0, len(c.Labels))
			for _, l := range c.Labels {
				parts = append(parts, l.Name)
			}
			fmt.Printf("  Labels: %s\n", strings.Join(parts, ", "))
		}
		if len(c.IDMembers) > 0 {
			fmt.Printf("  Members: %d assigned\n", len(c.IDMembers))
		}
		return nil
	}

	var action struct {
		Action *struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"action"`
	}
	if json.Unmarshal([]byte(text), &action) == nil && action.Action != nil {
		a := action.Action
		fmt.Printf("  Comment added: %s\n", a.Text)
		fmt.Printf("  ID: %s\n", a.ID)
		return nil
	}

	fmt.Println(text)
	return nil
}

func printChecklistsList(result *mcp.CallToolResult) error {
	var data struct {
		Checklists []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"checklists"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, ch := range data.Checklists {
		fmt.Printf("  %-40s %s\n", ch.Name, ch.ID)
	}
	return nil
}

func printActionsList(result *mcp.CallToolResult) error {
	var data struct {
		Actions []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"actions"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, a := range data.Actions {
		text := strings.ReplaceAll(a.Text, "\n", " ")
		if len(text) > 80 {
			text = text[:80] + "..."
		}
		fmt.Printf("  [%s] %s\n", a.Type, a.ID)
		if text != "" {
			fmt.Printf("    %s\n", text)
		}
	}
	return nil
}

func printListDetail(result *mcp.CallToolResult) error {
	var data struct {
		List struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"list"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	fmt.Printf("  %s\n", data.List.Name)
	fmt.Printf("  ID: %s\n", data.List.ID)
	return nil
}

func printMembersList(result *mcp.CallToolResult) error {
	var data struct {
		Members []struct {
			ID       string `json:"id"`
			FullName string `json:"full_name"`
			Username string `json:"username"`
		} `json:"members"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, m := range data.Members {
		fmt.Printf("  %-30s @%s  %s\n", m.FullName, m.Username, m.ID)
	}
	return nil
}

func printOrganizationsList(result *mcp.CallToolResult) error {
	var data struct {
		Organizations []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			URL         string `json:"url"`
		} `json:"organizations"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, o := range data.Organizations {
		fmt.Printf("  %-40s %s\n", o.DisplayName, o.ID)
		fmt.Printf("    URL: %s\n", o.URL)
	}
	return nil
}

func printLabelsList(result *mcp.CallToolResult) error {
	var data struct {
		Labels []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"labels"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, l := range data.Labels {
		name := l.Name
		if name == "" {
			name = "(no name)"
		}
		fmt.Printf("  %-30s %-10s %s\n", name, l.Color, l.ID)
	}
	return nil
}

func printCardLabelsList(result *mcp.CallToolResult) error {
	var data struct {
		Labels []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"labels"`
	}
	if err := json.Unmarshal([]byte(resultText(result)), &data); err != nil {
		fmt.Println(resultText(result))
		return nil
	}
	for _, l := range data.Labels {
		name := l.Name
		if name == "" {
			name = "(no name)"
		}
		fmt.Printf("  %-30s %-10s %s\n", name, l.Color, l.ID)
	}
	return nil
}

func resultText(result *mcp.CallToolResult) string {
	for _, c := range result.Content {
		if text, ok := c.(mcp.TextContent); ok {
			return text.Text
		}
	}
	return ""
}
