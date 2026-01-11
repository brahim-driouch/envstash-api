package repository

import (
	"context"
	"errors"
	"log"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContributorRepository struct {
	db *pgxpool.Pool
}

func NewContributorRepository(db *pgxpool.Pool) *ContributorRepository {
	return &ContributorRepository{
		db: db,
	}
}

//create contributor

func (r *ContributorRepository) CreateContributor(ctx context.Context, contributor models.CreateContributorRequest) error {
	cmdTag, err := r.db.Exec(ctx, queries.ContributorsQueries.CreateContributor, contributor.Fullname, contributor.Email, contributor.UserID)
	log.Println(contributor)
	if err != nil {
		log.Println(err)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("failed to inser new contributor")
	}
	return nil
}

func (r *ContributorRepository) GetContributorsByTeamIDs(ctx context.Context, teamIDs []string) (map[string][]models.Contributor, error) {
	if len(teamIDs) == 0 {
		return make(map[string][]models.Contributor), nil
	}
	for _, teamID := range teamIDs {
		// TODO: Implement database query using membersQueries.GetMembersByTeamIDs
		rows, err := r.db.Query(ctx, queries.ContributorsQueries.GetContributorsByTeamID, teamID)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	return make(map[string][]models.Contributor), nil
}

func (r *ContributorRepository) GetContributorByID(ctx context.Context, contributorID string) (*models.ContributorResponse, error) {
	var contributor models.ContributorResponse
	err := r.db.QueryRow(ctx, queries.ContributorsQueries.GetContributorByID, contributorID).Scan(
		&contributor.ID,
		&contributor.Fullname,
		&contributor.Email,
		&contributor.UserID,
		&contributor.IsVerified,
		&contributor.CreatedAt,
		&contributor.UpdatedAt,
		&contributor.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &contributor, nil
}

func (r *ContributorRepository) GetContributorsByUserID(ctx context.Context, userID string) ([]models.ContributorResponse, error) {
	rows, err := r.db.Query(ctx, queries.ContributorsQueries.GetContributorsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributors []models.ContributorResponse
	for rows.Next() {
		var contributor models.ContributorResponse
		err := rows.Scan(
			&contributor.ID,
			&contributor.Fullname,
			&contributor.Email,
			&contributor.UserID,
			&contributor.IsVerified,
			&contributor.CreatedAt,
			&contributor.UpdatedAt,
			&contributor.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		contributors = append(contributors, contributor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contributors, nil
}
