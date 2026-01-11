package models

import (
	"time"
)

// project
type Project struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required,min=1,max=255"`
	RepoURL     *string    `json:"repo_url,omitempty" db:"repo_url" validate:"omitempty,url"`
	Description string     `json:"description,omitempty" db:"description"`
	UserID      string     `json:"owner_id" db:"owner_id"`
	TeamID      *string    `json:"team_id,omitempty" db:"team_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// projects with stats
type ProjectWithStats struct {
	Project
	MemberCount int `json:"member_count" db:"member_count"`
	EnvVarCount int `json:"env_var_count" db:"env_var_count"`
}

// projects with details
type ProjectWithDetails struct {
	Project
	User        *User                 `json:"user,omitempty"`
	Team        *Team                 `json:"team,omitempty"`
	Members     []ProjectMember       `json:"members,omitempty"`
	EnvVars     []EnvironmentVariable `json:"env_vars,omitempty"`
	EnvVarCount int                   `json:"env_var_count"`
	MemberCount int                   `json:"member_count"`
}

// ProjectMember represents a user's membership in a project
type ProjectMember struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"project_id" db:"project_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	MemberID  string    `json:"member_id" db:"member_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ProjectMemberWithUser includes user details
type ProjectMemberWithUser struct {
	ProjectMember
	User *User `json:"user"`
}

// =====================================================
// REQUEST/RESPONSE STRUCTS
// =====================================================

// CreateProjectRequest represents the request to create a project
type CreateProjectRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	RepoURL     *string  `json:"repo_url,omitempty" validate:"omitempty,url"`
	Description string   `json:"description,omitempty" validate:"max=1000"`
	TeamID      *string  `json:"team_id,omitempty" validate:"omitempty,uuid"`
	Members     []string `json:"members,omitempty" validate:"dive,email"`
}

// UpdateProjectRequest represents the request to update a project
type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	RepoURL     *string `json:"repo_url,omitempty" validate:"omitempty,url"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	TeamID      *string `json:"team_id,omitempty" validate:"omitempty,uuid"`
}

// AddProjectMemberRequest represents the request to add a member to a project
type AddProjectMemberRequest struct {
	UserEmail string `json:"user_email" validate:"required,email"`
	Role      string `json:"role" validate:"required,oneof=editor viewer"`
}

// UpdateProjectMemberRequest represents the request to update a member's role
type UpdateProjectMemberRequest struct {
	Role string `json:"role" validate:"required,oneof=editor viewer"`
}

// ProjectResponse represents a project response
type ProjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	RepoURL     *string   `json:"repo_url,omitempty"`
	Description string    `json:"description,omitempty"`
	UserID      string    `json:"owner_id"`
	TeamID      *string   `json:"team_id,omitempty"`
	MemberCount int       `json:"member_count"`
	EnvVarCount int       `json:"env_var_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectListResponse represents a list of projects
type ProjectListResponse struct {
	Projects   []ProjectResponse `json:"projects"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// ProjectDetailsResponse represents detailed project information
type ProjectDetailsResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	RepoURL     *string                 `json:"repo_url,omitempty"`
	Description string                  `json:"description,omitempty"`
	User        UserResponse            `json:"user"`
	Team        *Team                   `json:"team,omitempty"`
	Members     []ProjectMemberResponse `json:"members"`
	EnvVars     []EnvVarResponse        `json:"env_vars"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// ProjectMemberResponse represents a project member in responses
type ProjectMemberResponse struct {
	UserID  string       `json:"user_id"`
	User    UserResponse `json:"user"`
	Role    string       `json:"role"`
	AddedAt time.Time    `json:"added_at"`
}

// =====================================================
// FILTER/QUERY STRUCTS
// =====================================================

// ProjectFilter represents filters for querying projects
type ProjectFilter struct {
	UserID         *string `json:"user_id,omitempty"`
	TeamID         *string `json:"team_id,omitempty"`
	Search         string  `json:"search,omitempty"` // Search in name/description
	IncludeDeleted bool    `json:"include_deleted"`
	Page           int     `json:"page" validate:"min=1"`
	PageSize       int     `json:"page_size" validate:"min=1,max=100"`
	SortBy         string  `json:"sort_by" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder      string  `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// =====================================================
// PERMISSION/ACCESS STRUCTS
// =====================================================

// ProjectAccess represents a user's access level to a project
type ProjectAccess struct {
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`   // owner, editor, viewer
	Source    string `json:"source"` // direct, team, owner
}

// ProjectPermissions defines what actions a user can perform
type ProjectPermissions struct {
	CanView          bool `json:"can_view"`
	CanEdit          bool `json:"can_edit"`
	CanDelete        bool `json:"can_delete"`
	CanManageMembers bool `json:"can_manage_members"`
	CanManageEnvVars bool `json:"can_manage_env_vars"`
}

// =====================================================
// HELPER STRUCTS
// =====================================================

// EnvironmentVariable (simplified - should be in env_vars package)
type EnvironmentVariable struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"project_id" db:"project_id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"` // Encrypted
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// EnvVarResponse for API responses (doesn't expose value)
type EnvVarResponse struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =====================================================
// CONSTANTS
// =====================================================

const (
	// Project member roles
	ProjectRoleOwner  = "owner"
	ProjectRoleEditor = "editor"
	ProjectRoleViewer = "viewer"
)

// IsValidProjectRole checks if a role is valid
func IsValidProjectRole(role string) bool {
	return role == ProjectRoleOwner ||
		role == ProjectRoleEditor ||
		role == ProjectRoleViewer
}

// GetProjectPermissions returns permissions based on role
func GetProjectPermissions(role string) ProjectPermissions {
	switch role {
	case ProjectRoleOwner:
		return ProjectPermissions{
			CanView:          true,
			CanEdit:          true,
			CanDelete:        true,
			CanManageMembers: true,
			CanManageEnvVars: true,
		}
	case ProjectRoleEditor:
		return ProjectPermissions{
			CanView:          true,
			CanEdit:          true,
			CanDelete:        false,
			CanManageMembers: false,
			CanManageEnvVars: true,
		}
	case ProjectRoleViewer:
		return ProjectPermissions{
			CanView:          true,
			CanEdit:          false,
			CanDelete:        false,
			CanManageMembers: false,
			CanManageEnvVars: false,
		}
	default:
		return ProjectPermissions{}
	}
}
