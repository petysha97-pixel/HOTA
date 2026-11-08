package models

type Candidate struct {
	ID           int    `json:"id"`
	Nickname     string `json:"nickname"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	MatchedStack int    `json:"matched_stack"`
	StackScore   int    `json:"stack_score"`
	RoleScore    int    `json:"role_score"`
	Match        int    `json:"match"`
	StackCount   int    `json:"-"`
}
