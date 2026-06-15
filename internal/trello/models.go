package trello

type Board struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Card struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	URL         string      `json:"url"`
	Desc        string      `json:"desc,omitempty"`
	Due         string      `json:"due,omitempty"`
	Start       string      `json:"start,omitempty"`
	DueComplete bool        `json:"dueComplete"`
	Closed      bool        `json:"closed"`
	IDBoard     string      `json:"idBoard,omitempty"`
	IDList      string      `json:"idList,omitempty"`
	IDMembers   []string    `json:"idMembers,omitempty"`
	Labels      []CardLabel `json:"labels,omitempty"`
}

type CardLabel struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard,omitempty"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Checklist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Action struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type Member struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
}

type BoardListResult struct {
	Boards []Board `json:"boards"`
}

type CardListResult struct {
	Cards []Card `json:"cards"`
}

type ListListResult struct {
	Lists []List `json:"lists"`
}

type GetBoardResult struct {
	Board Board `json:"board"`
}

type GetCardResult struct {
	Card Card `json:"card"`
}

type ChecklistListResult struct {
	Checklists []Checklist `json:"checklists"`
}

type ActionListResult struct {
	Actions []Action `json:"actions"`
}

type MemberListResult struct {
	Members []Member `json:"members"`
}

type CreateCardInput struct {
	ListID string
	Name   string
	Desc   string
	Due    string
}

type UpdateCardInput struct {
	Name string
	Desc string
	Due  string
}

type CreateCardResult struct {
	Card Card `json:"card"`
}

type UpdateCardResult struct {
	Card Card `json:"card"`
}

type MoveCardResult struct {
	Card Card `json:"card"`
}

type AddCommentResult struct {
	Action Action `json:"action"`
}

type BoardLabelListResult struct {
	Labels []Label `json:"labels"`
}

type CardLabelListResult struct {
	Labels []CardLabel `json:"labels"`
}

type LabelActionResult struct {
	Success bool `json:"success"`
}

type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
}

type SearchCardsResult struct {
	Cards []Card `json:"cards"`
}

type OrganizationListResult struct {
	Organizations []Organization `json:"organizations"`
}

type CreateBoardResult struct {
	Board Board `json:"board"`
}

type CreateListResult struct {
	List List `json:"list"`
}

type ArchiveCardResult struct {
	Card Card `json:"card"`
}

type UnarchiveCardResult struct {
	Card Card `json:"card"`
}
