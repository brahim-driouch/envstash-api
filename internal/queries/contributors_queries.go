package queries

var ContributorsQueries = struct {
	GetContributorsByTeamID  string
	CreateContributor        string
	GetContributorsByTeamIDs string
	GetContributorsByUserID  string
	GetContributorByID       string
}{
	GetContributorsByTeamID: `
		SELECT tm.member_id, tm.team_id, tm.added_at, tm.added_by, m.name as member_name
		FROM team_contributors tm
		JOIN contributors m ON tm.member_id = m.id
		WHERE tm.team_id = $1
		ORDER BY tm.added_at DESC`,
	GetContributorsByTeamIDs: `
		SELECT tm.member_id, tm.team_id, tm.added_at, tm.added_by, m.name as member_name
		FROM team_contributors tm
		JOIN contributors m ON tm.member_id = m.id
		WHERE tm.team_id = ANY($1)
		ORDER BY tm.added_at DESC
	`,
	CreateContributor: `
		INSERT INTO contributors (fullname, email, user_id)
		VALUES ($1, $2, $3)
	`,
	GetContributorsByUserID: `
		SELECT c.id, c.fullname, c.email, c.user_id, c.is_verified, c.created_at, c.updated_at, c.deleted_at
		FROM contributors c
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC
	`,
	GetContributorByID: `
		SELECT c.id, c.fullname, c.email, c.user_id, c.is_verified, c.created_at, c.updated_at, c.deleted_at
		FROM contributors c
		WHERE c.id = $1
	`,
}
