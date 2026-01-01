package queries

var ProjectQueries = struct {
	CreateProject string
}{
	CreateProject: `
		INSERT INTO projects (name, description, user_id)
		VALUES ($1, $2, $3)
	`,
}
