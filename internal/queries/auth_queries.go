package queries

var AuthQueries = struct {
	InsertUser           string
	UserExists           string
	FindUserByEmail      string
	FindUserByID         string
	CreateRefreshToken   string
	FindRefreshToken     string
	RevokeRefreshToken   string
	RevokeAllUserTokens  string
	DeleteExpiredTokens  string
	FindActiveUserTokens string
	DeleteUserToken      string
}{
	InsertUser: `
        INSERT INTO users ( fullname, email, password_hash)
        VALUES              ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `,
	FindUserByEmail: `
        SELECT id, fullname, email,  password_hash, is_verified, is_admin, created_at, updated_at
        FROM users
        WHERE email = $1
    `,
	FindUserByID: `
        SELECT id, email, fullname, password_hash, is_verified, is_admin, created_at, updated_at
        FROM users
        WHERE id = $1
    `,
	UserExists: `
        SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
    `,
	// Create a new refresh token
	CreateRefreshToken: `
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at, ip_address, user_agent) 
		VALUES ($1, $2, $3, $4, $5, $6 )
	`,

	// Find a valid (non-revoked, non-expired) refresh token
	FindRefreshToken: `
		SELECT id, user_id, token, expires_at, created_at, revoked_at, ip_address, user_agent 
		FROM refresh_tokens 
		WHERE token = $1 
		  AND revoked_at IS NULL 
		  AND expires_at > NOW()
	`,

	// Revoke a specific refresh token (logout)
	RevokeRefreshToken: `
		UPDATE refresh_tokens 
		SET revoked_at = NOW() 
		WHERE token = $1
	`,

	// Revoke all tokens for a user (logout from all devices)
	RevokeAllUserTokens: `
		UPDATE refresh_tokens 
		SET revoked_at = NOW() 
		WHERE user_id = $1 
		  AND revoked_at IS NULL
	`,

	// Delete expired and revoked tokens (cleanup job)
	DeleteExpiredTokens: `
		DELETE FROM refresh_tokens 
		WHERE expires_at < NOW() 
		   OR revoked_at IS NOT NULL
	`,

	// Find all active sessions for a user (session management)
	FindActiveUserTokens: `
		SELECT id, token, expires_at, created_at, ip_address, user_agent 
		FROM refresh_tokens 
		WHERE user_id = $1 
		  AND revoked_at IS NULL 
		  AND expires_at > NOW()
		ORDER BY created_at DESC
	`,

	// Delete a specific token by ID (for session management UI)
	DeleteUserToken: `
		UPDATE refresh_tokens 
		SET revoked_at = NOW() 
		WHERE id = $1 
		  AND user_id = $2
	`,
}
