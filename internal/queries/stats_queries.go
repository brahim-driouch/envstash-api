package queries

var Stats = struct {
	AllStats string
}{
	AllStats: `
	SELECT 
            (SELECT COUNT(*) FROM projects WHERE user_id = $1 AND deleted_at IS NULL) as projects_count,
            (SELECT COUNT(*) FROM teams WHERE user_id = $1 AND deleted_at IS NULL) as teams_count,
            (SELECT COUNT(DISTINCT tm.user_id) 
             FROM team_members tm 
             JOIN teams t ON tm.team_id = t.id 
             WHERE t.user_id = $1) as members_count,
            (SELECT COUNT(*) 
             FROM environment_variables ev 
             JOIN projects p ON ev.project_id = p.id 
             WHERE p.user_id = $1 AND p.deleted_at IS NULL) as env_vars_count
	`,
}
