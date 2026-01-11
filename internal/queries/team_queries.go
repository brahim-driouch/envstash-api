package queries

var TeamQueries = struct {
	CreateTeam       string
	GetTeamByID      string
	GetTeamsByUserID string
	UpdateTeam       string
	DeleteTeam       string
}{
	CreateTeam: `
		INSERT INTO teams (name, description, user_id, created_at)
		VALUES ($1, $2, $3, NOW())
		
	`,

	GetTeamByID: `
		SELECT id, name, user_id, created_at
		FROM teams
		WHERE id = $1
		AND deleted_at IS NULL
	`,

	GetTeamsByUserID: `
		SELECT id, name, user_id, created_at
		FROM teams
		WHERE user_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
	`,

	UpdateTeam: `
		UPDATE teams
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		AND deleted_at IS NULL
		RETURNING id, name, user_id, created_at
	`,

	DeleteTeam: `
		UPDATE teams
		SET deleted_at = NOW()
		WHERE id = $1
		`,
}
