package models

type UserStats struct {
	ProjectsCount int `json:"projects_count"`
	TeamsCount    int `json:"teams_count"`
	MembersCount  int `json:"members_count"`
	EnvVarsCount  int `json:"env_vars_count"`
}
