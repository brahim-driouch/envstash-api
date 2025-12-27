package queries

var UserQueries = struct {
	FindUserByID string
	UpdateUser   string
	DeleteUser   string
}{

	FindUserByID: `
        SELECT id, email, fullname, password_hash, is_verified, is_admin, created_at, updated_at
        FROM users
        WHERE id = $1
    `,

	UpdateUser: `
        UPDATE users
        SET fullname = $2, is_verified = $3, is_admin = $4, updated_at = NOW()
        WHERE id = $1
        RETURNING id, email, fullname, is_verified, is_admin, created_at, updated_at
    `,

	DeleteUser: `
        DELETE FROM users WHERE id = $1
    `,
}
