package queries

var ProjectQueries = struct {
	CreateProject        string
	GetAllProjects       string
	GetProjectByID       string
	GetProjectMembers    string
	GetProjectVars       string
	GetProjectsByTeamIDs string
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
	GetProjectByID: `
		SELECT id, name, description, user_id, created_at, updated_at
		FROM projects
		WHERE id = $1
	`,
	GetProjectMembers: `
		SELECT user_id, member_id, created_at
		FROM projects_members
		WHERE project_id = $1
	`,
	GetProjectVars: `
		SELECT id, project_id, key, value, created_at, updated_at, created_by
		FROM environment_variables
		WHERE project_id = $1
	`,
	GetProjectsByTeamIDs: `
		SELECT tp.project_id, tp.team_id, tp.added_at, tp.added_by, p.name as project_name
		FROM team_projects tp
		JOIN projects p ON tp.project_id = p.id
		WHERE tp.team_id = ANY($1)
		ORDER BY tp.added_at DESC
	`,
}
