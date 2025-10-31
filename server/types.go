package main

// TicketDialog represents the dialog data for ticket creation
type TicketDialog struct {
	TeamName    string `json:"team_name"`
	ProjectName string `json:"project_name"`
	Environment string `json:"environment"`
	Description string `json:"description"`
}
