package queries

var ProjectQueries = struct {
	CreateProject  string
	GetAllProjects string
}{
	CreateProject: `
		INSERT INTO projects (name, description, user_id)
		VALUES ($1, $2, $3)
	`,
	GetAllProjects: `
		SELECT id, name, description, user_id, created_at, updated_at
		FROM projects
		WHERE user_id = $1
	`,
}
