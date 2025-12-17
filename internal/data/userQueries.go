package data

var UserQueries = struct {
	InsertUser      string
	FindUserByID    string
	FindUserByEmail string
	UpdateUser      string
	DeleteUser      string
	UserExists      string
}{
	InsertUser: `
        INSERT INTO users ( fullname, email, password_hash)
        VALUES              ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `,

	FindUserByID: `
        SELECT id, email, fullname, password_hash, is_verified, is_admin, created_at, updated_at
        FROM users
        WHERE id = $1
    `,

	FindUserByEmail: `
        SELECT id, fullname, email,  password_hash, is_verified, is_admin, created_at, updated_at
        FROM users
        WHERE email = $1
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

	UserExists: `
        SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
    `,
}
