package models

type TeamWithMembersAndProjects struct {
	Team     Team          `json:"team"`
	Members  []Contributor `json:"members"`
	Projects []Project     `json:"projects"`
}
