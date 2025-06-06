package main

// Match holds information about a single football match including the teams and their goals.
type Match struct {
	Team1      string `json:"team1"`
	Team2      string `json:"team2"`
	Team1Goals string `json:"team1goals"`
	Team2Goals string `json:"team2goals"`
}

// APIResponse represents the structure of the response we get from the external football matches API.
type APIResponse struct {
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Data       []Match `json:"data"`
}

// TeamStats keeps track of a team's performance stats like wins, draws, losses, and goals.
type TeamStats struct {
	Team           string `json:"team"`
	MatchesPlayed  int    `json:"matches_played"`
	Wins           int    `json:"wins"`
	Draws          int    `json:"draws"`
	Losses         int    `json:"losses"`
	GoalsFor       int    `json:"goals_for"`
	GoalsAgainst   int    `json:"goals_against"`
	GoalDifference int    `json:"goal_difference"`
	Points         int    `json:"points"`
}
